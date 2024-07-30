package ddnsclient

import (
	"github.com/datahearth/ddnsclient/pkg/http"
	"github.com/datahearth/ddnsclient/pkg/providers/ovh"
	"github.com/sirupsen/logrus"
)

// Start create a new instance of ddns-client
func Start(logger logrus.FieldLogger) error {
	ddnsHTTP, err := http.NewHTTP(logger)
	if err != nil {
		return err
	}
	ovh, err := ovh.NewOVH(logger)
	if err != nil {
		return err
	}
	check := make(chan bool)
	c := make(chan bool)

	for {
		select {
		case <-check:

		case <-c:
			break
		}
	}

	return nil
}
