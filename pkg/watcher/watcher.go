package watcher

import (
	"time"

	"github.com/datahearth/ddnsclient/pkg/provider"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Watcher interface {
	Run(*time.Ticker, chan struct{}, chan error)
}

type watcher struct {
	logger       logrus.FieldLogger
	providers    []provider.Provider
	firstRun     bool
	webIP        string
	providerName string
}

// NewWatcher creates a watcher a given provider config
func NewWatcher(logger logrus.FieldLogger, w *utils.Watcher, webIP string) (Watcher, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if w == nil {
		return nil, utils.ErrNilWatcher
	}
	if webIP == "" {
		webIP = utils.DefaultURLs["webIP"]
	}

	providers := []provider.Provider{}
	for _, c := range w.Config {
		p, err := provider.NewProvider(logger, c, w.URL, w.Name)
		if err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	logger = logger.WithField("pkg", "watcher")

	return &watcher{
		logger:       logger,
		providers:    providers,
		webIP:        webIP,
		firstRun:     true,
		providerName: w.Name,
	}, nil
}

func (w *watcher) Run(t *time.Ticker, chClose chan struct{}, chErr chan error) {
	logger := w.logger.WithField("component", "Run")

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
			logger.WithField("provider", w.providerName).Infoln("Close watcher channel triggered. Ticker stopped")
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

	logger.Debugln("Checking server IP...")
	srvIP, err := utils.RetrieveServerIP(w.webIP)
	if err != nil {
		return err
	}

	for _, p := range w.providers {
		if err := p.UpdateSubdomains(srvIP); err != nil {
			return err
		}
	}

	logger.Infof("[%s] DDNS check finished\n", w.providerName)
	return nil
}
