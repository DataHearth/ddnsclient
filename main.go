package ddnsclient

import (
	"errors"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/providers/google"
	"github.com/datahearth/ddnsclient/pkg/providers/ovh"
	"github.com/datahearth/ddnsclient/pkg/watcher"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
)

var (
	ErrSbsLen             = errors.New("subdomains len is 0")
	ErrInvalidProvider    = errors.New("invalid provider name")
	ErrWatchersConfigLen  = errors.New("watcher configuration needs at least one [watcherd] and [providerd] configuration")
	ErrWatcherCreationLen = errors.New("no valid watchers were created. Checkout [watchers] configuration and its [providers] configuration")
)

// Start create a new instance of ddns-client
func Start(logger logrus.FieldLogger, config ClientConfig) error {
	log := logger.WithFields(logrus.Fields{
		"pkg":       "ddnsclient",
		"component": "root",
	})

	fields := reflect.ValueOf(config.Watchers)
	ws := []watcher.Watcher{}

	// * check providers and watchers config
	// todo: invalid condition but while still exit in the next step. To be corrected
	if fields.NumField() == 0 || reflect.ValueOf(config.Providers).NumField() == 0 {
		return ErrWatchersConfigLen
	}

	for i := 0; i < fields.NumField(); i++ {
		providerName := strings.ToLower(fields.Type().Field(i).Name)

		w, err := CreateWatcher(providerName, config.WebIP, logger, config.Watchers, config.Providers, config.PendingDnsPropagation)
		if err != nil {
			logger.Warnf("Provider error: %v. Skipping...\n", err.Error())
			continue
		}

		ws = append(ws, w)
	}

	// * check for valid created watchers
	if len(ws) == 0 {
		return ErrWatcherCreationLen
	}

	// * create signal watcher
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer close(sigc)

	// * create close and error channel
	chClose := make(chan struct{})
	chErr := make(chan error)
	defer close(chClose)
	defer close(chErr)

	log.Infoln("Start watching periodically for changes!")
	// * run every created watchers in goroutines
	for _, w := range ws {
		tickTime := config.UpdateTime
		if tickTime == 0 {
			tickTime = 180
		}

		t := time.NewTicker(time.Duration(tickTime) * time.Second)
		go w.Run(t, chClose, chErr)
	}

	// * listening for errors and exit signal
	for {
		select {
		case err := <-chErr:
			log.Errorln(err.Error())
			continue
		case <-sigc:
			log.Infoln("Interrupt signal received. Stopping watcher...")
			chClose <- struct{}{}
			return nil
		}
	}
}

func CreateWatcher(provider, webIP string, logger logrus.FieldLogger, wc WatcherConfig, ps Providers, pendingDnsPropagation int) (watcher.Watcher, error) {
	var sbs []string
	var p providers.Provider
	var err error

	// * check for implemented providers
	switch provider {
	case "ovh":
		log.Debugln("create OVH provider")
		p, err = ovh.NewOVH(logger, &ps.Ovh)
		if err != nil {
			return nil, err
		}
		sbs = wc.Ovh

	case "google":
		log.Debugln("create GOOGLE provider")
		p, err = google.NewGoogle(logger, &ps.Google)
		if err != nil {
			return nil, err
		}
		sbs = wc.Google

	default:
		return nil, ErrInvalidProvider
	}

	if len(sbs) == 0 {
		return nil, ErrSbsLen
	}

	// * create provider's watcher
	w, err := watcher.NewWatcher(logger, p, sbs, webIP, provider, pendingDnsPropagation)
	if err != nil {
		return nil, err
	}

	return w, nil
}
