package google

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	// ErrNilGoogleConfig is thrown when GOOGLE configuration is empty
	ErrNilGoogleConfig = errors.New("GOOGLE config is mandatory")
	// ErrInvalidConfig is thrown when no username and password are provided and URL doesn't contains them
	ErrInvalidConfig = errors.New("username and password are required if url doesn't contains them")
	// ErrReadBody is thrown when reader failed to read response body
	ErrReadBody = errors.New("failed to read response body")
)

// GoogleConfig is the struct for the yaml configuration file
type GoogleConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
}

type google struct {
	config *GoogleConfig
	logger logrus.FieldLogger
}

// NewGoogle returns a new instance of the GOOGLE provider
func NewGoogle(logger logrus.FieldLogger, googleConfig *GoogleConfig) (providers.Provider, error) {
	if googleConfig == nil {
		return nil, ErrNilGoogleConfig
	}
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if (googleConfig.Username == "" && googleConfig.Password == "") && !strings.Contains(googleConfig.URL, "@") {
		return nil, ErrInvalidConfig
	}

	logger = logger.WithField("pkg", "provider-google")

	return &google{
		config: googleConfig,
		logger: logger,
	}, nil
}

// UpdateIP updates the subdomain A record
func (g *google) UpdateIP(subdomain, ip string) error {
	newURL := strings.ReplaceAll(g.config.URL, "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)
	logger := g.logger.WithFields(logrus.Fields{
		"component":      "update-ip",
		"ovh-update-url": newURL,
		"subdomain":      subdomain,
		"new-ip":         ip,
	})

	// * create POST request
	req, err := http.NewRequest("POST", newURL, nil)
	if err != nil {
		return utils.ErrCreateNewRequest
	}
	// * use basic auth if config is set
	if g.config.Username != "" && g.config.Password != "" {
		logger.Debugln("username and password passed in config. Use basic auth")
		req.SetBasicAuth(g.config.Username, g.config.Password)
	}

	// * perform POST request
	logger.Debugln("calling Google DynDNS to update subdomain IP")
	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return utils.ErrUpdateRequest
	}

	// * read response body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadBody
	}

	// * check for error response
	// doc: https://support.google.com/domains/answer/6147083?hl=en#zippy=%2Cusing-the-api-to-update-your-dynamic-dns-record
	// todo: check why the hell do I need to use () for conditions here !!!!
	if (strings.Contains(string(b), "good "+ip) != true) || (strings.Contains(string(b), "nochg "+ip) != false) {
		return errors.New("failed to update subdomain ip. Google error: " + string(b))
	}

	return nil
}
