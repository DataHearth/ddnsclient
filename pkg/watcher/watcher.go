package watcher

import (
	"fmt"
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
	logger            logrus.FieldLogger
	provider          providers.Provider
	subdomains        []subdomain.Subdomain
	domain            string
	webIP             string
	firstRun          bool
	pendingSubdomains subdomain.PendingSubdomains
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
	var sbs []string
	if sb, ok := viper.GetStringMap("watcher")["subdomains"].([]interface{}); ok {
		for _, v := range sb {
			sbs = append(sbs, fmt.Sprint(v))
		}
	}

	sbs = utils.AggregateSubdomains(sbs, domain)
	subdomains := make([]subdomain.Subdomain, len(sbs))
	for i, sd := range sbs {
		sub, err := subdomain.NewSubdomain(logger, sd)
		if err != nil {
			return nil, err
		}

		subdomains[i] = sub
	}

	return &watcher{
		logger:            logger,
		provider:          provider,
		domain:            domain,
		subdomains:        subdomains,
		webIP:             webIP,
		firstRun:          true,
		pendingSubdomains: make(map[time.Time]subdomain.Subdomain),
	}, nil
}

func (w *watcher) Run(t *time.Ticker, chClose chan struct{}, chErr chan error) {
	logger := w.logger.WithField("component", "Run")

	go w.checkPendingSubdomains(chClose)
	if w.firstRun {
		if err := w.runDDNSCheck(); err != nil {
			chErr <- err
		}
		w.firstRun = false
	}

	for {
		select {
		case <-chClose:
			t.Stop()
			logger.Infoln("Close watcher channel triggered. Ticker stoped")
			return
		case <-t.C:
			if err := w.runDDNSCheck(); err != nil {
				chErr <- err
			}
		}
	}
}

func (w *watcher) runDDNSCheck() error {
	logger := w.logger.WithField("component", "runDDNSCheck")

	logger.Infoln("Starting DDNS check...")

	srvIP, err := utils.RetrieveServerIP(w.webIP)
	if err != nil {
		return err
	}
	logger.Debugln("Checking server IP...")

	srvIP = "109.14.53.74" // tmp

	for _, sb := range w.subdomains {
		if sb.SubIsPending(w.pendingSubdomains) {
			continue
		}

		logger.Debugf("Checking subdomain %s...\n", sb.GetSubdomainAddr())
		ok, err := sb.CheckIPAddr(srvIP)
		if err != nil {
			return err
		}
		subAddr := sb.GetSubdomainAddr()
		if !ok {
			logger.WithFields(logrus.Fields{
				"server-ip":         srvIP,
				"subdomain-address": subAddr,
			}).Infoln("IP addresses doesn't match. Updating subdomain's ip...")
			if err := w.provider.UpdateIP(subAddr, srvIP); err != nil {
				logger.WithError(err).WithFields(logrus.Fields{
					"server-ip":         srvIP,
					"subdomain-address": subAddr,
				}).Errorln("failed to update subdomain's ip")
				return err
			}
			logger.WithFields(logrus.Fields{
				"server-ip":         srvIP,
				"subdomain-address": subAddr,
			}).Infoln("Subdomain's ip updated! Removing from checks for 5 mins")

			w.pendingSubdomains[time.Now()] = sb

			continue
		}
		logger.Debugf("%s is up to date. \n", subAddr)
	}

	logger.Infoln("DDNS check finished")
	return nil
}

func (w *watcher) checkPendingSubdomains(chClose chan struct{}) {
	logger := w.logger.WithField("component", "checkPendingSubdomains")
	t := time.NewTicker(time.Second * 10)

	logger.Debugln("Start checking for pending subdomains...")
	for {
		select {
		case <-chClose:
			logger.Debugln("Close pending subdomains")
			return
		case <-t.C:
			logger.Debugln("Checking pending subdomains...")
			if delSbs := subdomain.CheckPendingSubdomains(w.pendingSubdomains, time.Now()); delSbs != nil {
				w.pendingSubdomains = subdomain.DeletePendingSubdomains(delSbs, w.pendingSubdomains)
				logger.Debugln("Pendings subdomains found. Cleaned.")
			}
		}
	}
}
