package ddnsclient

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/datahearth/ddnsclient/pkg/providers/ovh"
	"github.com/datahearth/ddnsclient/pkg/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Start create a new instance of ddns-client
func Start(logger logrus.FieldLogger, config ClientConfig) error {
	log := logger.WithFields(logrus.Fields{
		"pkg":       "ddnsclient",
		"component": "root",
	})

	log.Debugln("create OVH provider")
	ovh, err := ovh.NewOVH(logger)
	if err != nil {
		return err
	}

	log.Debugln("creating watcher with OVH provider")
	w, err := watcher.NewWatcher(logger, ovh, viper.GetString("web-ip"))
	if err != nil {
		return err
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	chClose := make(chan struct{})
	chErr := make(chan error)
	defer close(chClose)
	defer close(chErr)
	defer close(sigc)

	log.Infoln("Start watching periodically for changes!")
	t := time.NewTicker(viper.GetDuration("update-time") * time.Second)
	go w.Run(t, chClose, chErr)

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
