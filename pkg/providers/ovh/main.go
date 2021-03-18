package ovh

import (
	"net/http"
	"strings"

	"github.com/datahearth/ddnsclient"
	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

type ovh struct {
	ovhConfig *ddnsclient.Ovh
	logger    logrus.FieldLogger
}

// NewOVH returns a new instance of the OVH provider
func NewOVH(logger logrus.FieldLogger, ovhConfig *ddnsclient.Ovh) (providers.Provider, error) {
	if ovhConfig == nil {
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
	newURL := strings.ReplaceAll(ovh.ovhConfig.URL, "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)
	logger := ovh.logger.WithFields(logrus.Fields{
		"component":      "update-ip",
		"ovh-update-url": newURL,
	})

	// * create GET request
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return utils.ErrCreateNewRequest
	}
	req.SetBasicAuth(ovh.ovhConfig.Username, ovh.ovhConfig.Password)

	// * perform GET request
	logger.WithFields(logrus.Fields{
		"subdomain": subdomain,
		"new-ip":    ip,
	}).Debugln("calling OVH DynHost to update subdomain IP")
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
