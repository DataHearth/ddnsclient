package ovh

import (
	"net/http"
	"strings"

	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ovh struct {
	ovhConfig config
	logger    logrus.FieldLogger
}

// NewOVH returns a new instance of the OVH provider
func NewOVH(logger logrus.FieldLogger) (providers.Provider, error) {
	var ovhConfig config
	if c, ok := viper.GetStringMap("provider")["ovh"].(config); ok {
		ovhConfig = c
	} else {
		return nil, ErrNilOvhConfig
	}

	if logger == nil {
		return nil, utils.ErrNilLogger
	}

	return &ovh{
		ovhConfig: ovhConfig,
		logger:    logger,
	}, nil
}

func (ovh *ovh) UpdateIP(subdomain, ip string) error {
	newURL := strings.ReplaceAll(ovh.ovhConfig["url"].(string), "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)

	// * create GET request
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(ovh.ovhConfig["username"].(string), ovh.ovhConfig["password"].(string))

	// * perform GET request
	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return utils.ErrWrongStatusCode
	}

	return nil
}
