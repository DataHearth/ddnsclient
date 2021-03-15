package watcher

import (
	"time"

	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/subdomain"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Watcher interface {
	Run(*time.Ticker, chan struct{}, chan error)
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
		return nil, utils.ErrNilProvider
	}
	if webIP == "" {
		webIP = "http://dynamicdns.park-your-domain.com/getip"
	}
	logger = logger.WithField("pkg", "watcher")

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

func (w *watcher) Run(t *time.Ticker, chClose chan struct{}, chErr chan error) {
	logger := w.logger.WithField("component", "run")

	for {
		select {
		case <-chClose:
			t.Stop()
			logger.Infoln("Close watcher channel triggered. Ticker stoped")
			return
		case <-t.C:
			logger.Infoln("Starting DDNS check")
			srvIP, err := utils.RetrieveServerIP(w.webIP)
			if err != nil {
				chErr <- err
				continue
			}

			logger.WithField("server-ip", srvIP).Debugln("Server IP retrieved. Checking subdomains...")
			for _, sd := range w.subdomains {
				ok, err := sd.CheckIPAddr(srvIP)
				if err != nil {
					logger.WithError(err).WithField("server-ip", srvIP).Errorln("failed to check ip addresses")
					chErr <- err
					continue
				}
				if !ok {
					subIP := sd.GetSubdomainIP()
					logger.WithFields(logrus.Fields{
						"server-ip":    srvIP,
						"subdomain-ip": subIP,
					}).Infoln("IP addresses doesn't match. Updating subdomain's ip...")
					if err := w.provider.UpdateIP(subIP, srvIP); err != nil {
						logger.WithError(err).WithFields(logrus.Fields{
							"server-ip":    srvIP,
							"subdomain-ip": subIP,
						}).Errorln("failed to update subdomain's ip")
						chErr <- err
						continue
					}
					logger.WithFields(logrus.Fields{
						"server-ip":    srvIP,
						"subdomain-ip": subIP,
					}).Infoln("Subdomain updated successfully!")
				}
			}
		}
	}
}
