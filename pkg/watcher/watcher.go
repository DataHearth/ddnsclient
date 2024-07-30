package watcher

import (
	"time"

	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/subdomain"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Watcher interface {
	Run(*time.Ticker, chan struct{}, chan error)
}

type watcher struct {
	logger                logrus.FieldLogger
	provider              providers.Provider
	subdomains            []subdomain.Subdomain
	pendingSubdomains     subdomain.PendingSubdomains
	firstRun              bool
	pendingDnsPropagation int
	webIP                 string
	providerName          string
}

// NewWatcher creates a watcher a given provider and its subdomains
func NewWatcher(logger logrus.FieldLogger, provider providers.Provider, sbs []string, webIP, providerName string, pendingDnsPropagation int) (Watcher, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if provider == nil {
		return nil, utils.ErrNilProvider
	}
	if webIP == "" {
		webIP = "http://dynamicdns.park-your-domain.com/getip"
	}
	if pendingDnsPropagation == 0 {
		pendingDnsPropagation = 180
	}
	logger = logger.WithField("pkg", "watcher")

	subdomains := make([]subdomain.Subdomain, len(sbs))
	for i, sb := range sbs {
		sub, err := subdomain.NewSubdomain(logger, sb)
		if err != nil {
			return nil, err
		}

		subdomains[i] = sub
	}

	return &watcher{
		logger:                logger,
		provider:              provider,
		subdomains:            subdomains,
		webIP:                 webIP,
		firstRun:              true,
		pendingSubdomains:     make(map[time.Time]subdomain.Subdomain),
		pendingDnsPropagation: pendingDnsPropagation,
		providerName:          providerName,
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
			logger.Infoln("Close watcher channel triggered. Ticker stopped")
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

	logger.Infof("Starting [%s] DDNS check...\n", w.providerName)

	srvIP, err := utils.RetrieveServerIP(w.webIP)
	if err != nil {
		return err
	}
	logger.Debugln("Checking server IP...")

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

	logger.Infof("[%s] DDNS check finished\n", w.providerName)
	return nil
}

func (w *watcher) checkPendingSubdomains(chClose chan struct{}) {
	logger := w.logger.WithField("component", "checkPendingSubdomains")
	t := time.NewTicker(time.Second * time.Duration(w.pendingDnsPropagation))

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
