package ddnsclient

import (
	"time"

	"github.com/datahearth/ddnsclient/pkg/providers/ovh"
	"github.com/datahearth/ddnsclient/pkg/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Start create a new instance of ddns-client
func Start(logger logrus.FieldLogger) error {
	ovh, err := ovh.NewOVH(logger)
	if err != nil {
		return err
	}

	w, err := watcher.NewWatcher(logger, ovh, viper.GetString("web-ip"))
	if err != nil {
		return err
	}

	c := make(chan bool)
	go w.Run(c)
	for {
		
	}

	return nil
}
