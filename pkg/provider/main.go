package provider

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Provider interface {
	UpdateSubdomains(ip string) error
	updateSubdomain(subdomain, ip string) error
	retrieveSubdomainIP(addr string) (string, error)
}

type provider struct {
	logger logrus.FieldLogger
	config utils.Config
	name   string
	url    string
}

// NewProvider creates a new instance of the `Provider` interface
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

// UpdateSubdomains will watch for every defined subdomains if they need an update.
// If so, it'll trigger an update
func (p *provider) UpdateSubdomains(srvIP string) error {
	for _, sb := range p.config.Subdomains {
		ip, err := p.retrieveSubdomainIP(sb)
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
			p.logger.WithError(err).WithFields(logrus.Fields{
				"component": "UpdateSubdomains",
				"subdomain": sb,
				"new-ip":    srvIP,
			}).Errorln("failed to update subdomain ip")
		}
	}

	return nil
}

func (p *provider) updateSubdomain(subdomain, ip string) error {
	tokenBased := p.config.Token != "" && (p.config.Username == "" && p.config.Password == "")

	newURL := strings.ReplaceAll(p.url, "SUBDOMAIN", subdomain)
	newURL = strings.ReplaceAll(newURL, "NEWIP", ip)
	if tokenBased {
		newURL = strings.ReplaceAll(newURL, "TOKEN", p.config.Token)
	}
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
	if !tokenBased {
		req.SetBasicAuth(p.config.Username, p.config.Password)
	}

	logger.Debugln("calling DDNS provider for subdomain update")
	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.ContentLength != 0 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := p.checkResponse(b, tokenBased, ip); err != nil {
			return err
		}
	}

	if resp.StatusCode != 200 {
		return utils.ErrWrongStatusCode
	}

	return nil
}

func (p *provider) retrieveSubdomainIP(addr string) (string, error) {
	ips, err := net.LookupIP(addr)
	if err != nil {
		return "", err
	}

	if len(ips) != 1 {
		return "", utils.ErrIpLenght
	}

	ip := ips[0].String()
	if strings.Contains(ip, ":") {
		ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			return "", utils.ErrSplitAddr
		}
	}

	return ip, nil
}

func (p *provider) checkResponse(body []byte, tokenBased bool, ip string) error {
	var invalidResponse error

	if tokenBased {
		if !strings.Contains(string(body), "OK") {
			if strings.Contains(string(body), "KO") {
				invalidResponse = fmt.Errorf("invalid body response.\n Body response: %v", string(body))
			} else {
				invalidResponse = fmt.Errorf("unknown body response. Please fill a issue if you think this is an error.\n Body response: %v", string(body))
			}
		}
	} else {
		if !strings.Contains(string(body), "good "+ip) && !strings.Contains(string(body), "nochg "+ip) {
			invalidResponse = fmt.Errorf("invalid body response.\n Body response: %v", string(body))
		}
	}

	return invalidResponse
}
