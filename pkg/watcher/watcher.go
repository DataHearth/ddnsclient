package watcher

import (
	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/subdomain"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Watcher interface {
	Run(chan bool) error
}

type watcher struct {
	logger     logrus.FieldLogger
	provider   providers.Provider
	subdomains []subdomain.Subdomain
	domain     string
	webIP      string
}

func NewWatcher(logger logrus.FieldLogger, provider providers.Provider, webIP string) (Watcher, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if provider == nil {
		return nil, ErrNilProvider
	}
	if webIP == "" {
		webIP = "http://dynamicdns.park-your-domain.com/getip"
	}
	domain := viper.GetStringMap("watcher")["domain"].(string)
	sbs := utils.AggregateSubdomains(viper.GetStringMap("watcher")["subdomains"].([]string), domain)
	subdomains := make([]subdomain.Subdomain, len(sbs))
	for _, sd := range sbs {
		sub, err := subdomain.NewSubdomain(logger, sd)
		if err != nil {
			return nil, err
		}

		subdomains = append(subdomains, sub)
	}

	return &watcher{
		logger:     logger,
		provider:   provider,
		domain:     domain,
		subdomains: subdomains,
		webIP:      webIP,
	}, nil
}

func (w *watcher) Run(close chan bool) error {
	for {
		select {
		case <-close:
			return nil
		default:
			srvIP, err := utils.RetrieveServerIP(w.webIP)
			if err != nil {
				return err
			}

			for _, sd := range w.subdomains {
				ok, err := sd.CheckIPAddr(srvIP)
				if err != nil {
					return err
				}
				if !ok {
					w.provider.UpdateIP(sd.GetSubdomainIP(), srvIP)
				}
			}
		}
	}
}
