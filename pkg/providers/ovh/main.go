package ovh

import (
	"errors"
	"net/http"
	"strings"

	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

// ErrNilOvhConfig is thrown when OVH configuration is empty
var ErrNilOvhConfig = errors.New("OVH config is mandatory")

type OvhConfig struct {
	URL      string `mapstructure:"url,omitempty"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ovh struct {
	config *OvhConfig
	logger logrus.FieldLogger
}

// NewOVH returns a new instance of the OVH provider
func NewOVH(logger logrus.FieldLogger, ovhConfig *OvhConfig) (providers.Provider, error) {
	if ovhConfig == nil {
		return nil, ErrNilOvhConfig
	}
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if ovhConfig.URL == "" {
		ovhConfig.URL = "http://www.ovh.com/nic/update?system=dyndns&hostname=SUBDOMAIN&myip=NEWIP"
	}

	logger = logger.WithField("pkg", "provider-ovh")

	return &ovh{
		config: ovhConfig,
		logger: logger,
	}, nil
}

func (ovh *ovh) UpdateIP(subdomain, ip string) error {
	newURL := strings.ReplaceAll(ovh.config.URL, "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)
	logger := ovh.logger.WithFields(logrus.Fields{
		"component":      "update-ip",
		"ovh-update-url": newURL,
		"subdomain":      subdomain,
		"new-ip":         ip,
	})

	// * create GET request
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return utils.ErrCreateNewRequest
	}
	req.SetBasicAuth(ovh.config.Username, ovh.config.Password)

	// * perform GET request
	logger.Debugln("calling OVH DynHost to update subdomain IP")
	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return utils.ErrUpdateRequest
	}
	if resp.StatusCode != 200 {
		return utils.ErrWrongStatusCode
	}

	return nil
}
