package provider

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

// Provider is the default interface for all providers
type Provider interface {
	UpdateSubdomains(ip string) error
}

type provider struct {
	logger logrus.FieldLogger
	config utils.Config
	name   string
	url    string
}

func NewProvider(logger logrus.FieldLogger, config utils.Config, url, name string) (Provider, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if name == "" {
		return nil, utils.ErrInvalidName
	}
	if url == "" {
		if utils.DefaultURLs[name] == "" {
			return nil, utils.ErrInvalidURL
		}
		url = utils.DefaultURLs[name]
	}
	logger = logger.WithField("pkg", "providers")

	return &provider{
		config: config,
		logger: logger,
		name:   name,
		url:    url,
	}, nil
}

func (p *provider) UpdateSubdomains(srvIP string) error {
	for _, sb := range p.config.Subdomains {
		ip, err := utils.RetrieveSubdomainIP(sb)
		if err != nil {
			return err
		}

		if ip == srvIP {
			continue
		}

		p.logger.WithFields(logrus.Fields{
			"component":         "UpdateSubdomains",
			"server-ip":         srvIP,
			"subdomain-address": ip,
			"subdomain":         sb,
		}).Infoln("IP addresses doesn't match. Updating subdomain's ip...")
		if err := p.updateSubdomain(sb, srvIP); err != nil {
			if err != utils.ErrReadBody && err != utils.ErrWrongStatusCode {
				return err
			}
			p.logger.WithError(err).WithFields(logrus.Fields{
				"component": "UpdateSubdomains",
				"subdomain": sb,
				"new-ip":    srvIP,
			}).Warnln("failed to update subdomain ip")
		}
	}

	return nil
}

func (p *provider) updateSubdomain(subdomain, ip string) error {
	newURL := strings.ReplaceAll(p.url, "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)
	logger := p.logger.WithFields(logrus.Fields{
		"component":   "UpdateIP",
		"updated-url": newURL,
		"subdomain":   subdomain,
		"new-ip":      ip,
	})

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return utils.ErrCreateNewRequest
	}
	req.SetBasicAuth(p.config.Username, p.config.Password)

	logger.Debugln("calling DDNS provider for subdomain update")
	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return utils.ErrUpdateRequest
	}

	if resp.ContentLength != 0 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return utils.ErrReadBody
		}

		if !strings.Contains(string(b), "good "+ip) && !strings.Contains(string(b), "nochg "+ip) {
			return errors.New("failed to update subdomain ip. Error: " + string(b))
		}
	}

	if resp.StatusCode != 200 {
		return utils.ErrWrongStatusCode
	}

	return nil
}
