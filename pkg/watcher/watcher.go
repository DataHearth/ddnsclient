package watcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/datahearth/ddnsclient/pkg/provider"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Watcher interface {
	Run(*time.Ticker, chan struct{}, chan error)
	runDDNSCheck() error
	retrieveServerIP() (string, error)
}

type watcher struct {
	logger       logrus.FieldLogger
	providers    []provider.Provider
	firstRun     bool
	webIP        string
	providerName string
}

// NewWatcher creates a new instance of the `Watcher` interface
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

// Run will trigger a subdomain update every X seconds (defined in config `update-time`)
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
	srvIP, err := w.retrieveServerIP()
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

func (w *watcher) retrieveServerIP() (string, error) {
	rsp, err := http.Get(w.webIP)
	if err != nil {
		return "", fmt.Errorf("%v: %v", utils.ErrGetServerIP, err)
	}
	if rsp.StatusCode != 200 {
		return "", fmt.Errorf("%v: %v", utils.ErrWrongStatusCode, rsp.Status)
	}

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %v", utils.ErrParseHTTPBody, err)
	}

	var ip string
	if rsp.Header.Get("content-type") == "application/json" {
		body := make(map[string]interface{})
		if err := json.Unmarshal(b, &body); err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON body: %v", err)
		}

		if v, ok := body["ip"]; ok {
			if v, ok := v.(string); ok {
				ip = v
			} else {
				return "", fmt.Errorf("\"ip\" field found in JSON response but is not a string value")
			}
		} else {
			for _, v := range body {
				if v, ok := v.(string); ok {
					if parsedIP := net.ParseIP(v); parsedIP != nil {
						ip = parsedIP.String()
					}
				}
			}

			if ip == "" {
				return "", fmt.Errorf("no field found with a valid IPv4 address in JSON response")
			}
		}
	} else {
		body := string(b)
		if strings.Count(body, "\n") > 0 {
			for _, l := range strings.Split(body, "\n") {
				if parsedIP := net.ParseIP(l); parsedIP != nil {
					ip = parsedIP.String()
					break
				}
			}

			if ip == "" {
				return "", fmt.Errorf("no field found with a valid IPv4 address in body response")
			}
		} else {
			if body == "" {
				return "", fmt.Errorf("body is empty")
			}
			ip = body
		}

	}

	return ip, nil
}
