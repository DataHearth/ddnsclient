package ddnsclient

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/datahearth/ddnsclient/pkg/watcher"
	"github.com/sirupsen/logrus"
)

var (
	ErrSbsLen             = errors.New("subdomains len is 0")
	ErrInvalidProvider    = errors.New("invalid provider name")
	ErrWatchersConfigLen  = errors.New("watcher configuration needs at least one [watcherd] and [providerd] configuration")
	ErrWatcherCreationLen = errors.New("no valid watchers were created. Checkout [watchers] configuration and its [providers] configuration")
)

// Start create a new instance of ddns-client
func Start(logger logrus.FieldLogger, config utils.ClientConfig) error {
	log := logger.WithFields(logrus.Fields{
		"pkg":       "ddnsclient",
		"component": "root",
	})

	ws := make([]watcher.Watcher, 0, len(config.Watchers))
	for _, cw := range config.Watchers {
		w, err := watcher.NewWatcher(logger, &cw, config.WebIP)
		if err != nil {
			logger.Warnf("Provider error: %v. Skipping...\n", err.Error())
			continue
		}

		ws = append(ws, w)
	}

	if len(ws) == 0 {
		return ErrWatcherCreationLen
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer close(sigc)

	chClose := make(chan struct{})
	chErr := make(chan error)
	defer close(chClose)
	defer close(chErr)

	logger.Infoln("Start watching periodically for changes!")
	for _, w := range ws {
		tickTime := config.UpdateTime
		if tickTime == 0 {
			tickTime = 180
		}

		t := time.NewTicker(time.Duration(tickTime) * time.Second)
		go w.Run(t, chClose, chErr)
	}

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
