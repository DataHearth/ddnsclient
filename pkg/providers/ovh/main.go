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
	ovhConfig utils.ProviderConfig
	logger    logrus.FieldLogger
}

// NewOVH returns a new instance of the OVH provider
func NewOVH(logger logrus.FieldLogger) (providers.Provider, error) {
	var ovhConfig utils.ProviderConfig
	if c, ok := viper.GetStringMap("provider")["ovh"].(utils.ProviderConfig); ok {
		ovhConfig = c
	} else {
		return nil, utils.ErrNilOvhConfig
	}
	if logger == nil {
		return nil, utils.ErrNilLogger
	}

	logger = logger.WithField("pkg", "provider-ovh")

	return &ovh{
		ovhConfig: ovhConfig,
		logger:    logger,
	}, nil
}

func (ovh *ovh) UpdateIP(subdomain, ip string) error {
	newURL := strings.ReplaceAll(ovh.ovhConfig["url"].(string), "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)
	logger := ovh.logger.WithFields(logrus.Fields{
		"component":      "update-ip",
		"ovh-update-url": newURL,
	})

	// * create GET request
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		logger.WithError(err).WithField("request-type", "GET").Errorln(utils.ErrCreateNewRequest.Error())
		return utils.ErrCreateNewRequest
	}
	req.SetBasicAuth(ovh.ovhConfig["username"].(string), ovh.ovhConfig["password"].(string))

	// * perform GET request
	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		logger.WithError(err).Errorln(utils.ErrUpdateRequest.Error())
		return utils.ErrUpdateRequest
	}
	if resp.StatusCode != 200 {
		logger.WithField("status-code", resp.StatusCode).Errorln(utils.ErrWrongStatusCode.Error())
		return utils.ErrWrongStatusCode
	}

	return nil
}
