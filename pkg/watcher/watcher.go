package watcher

import (
	"github.com/datahearth/ddnsclient/pkg/http"
	"github.com/datahearth/ddnsclient/pkg/providers"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Watcher interface {
	Run() error
}

type watcher struct {
	logger     logrus.FieldLogger
	provider   providers.Provider
	http       http.HTTP
	subdomains []string
	domain     string
}

func NewWatcher(logger logrus.FieldLogger, provider providers.Provider, http http.HTTP) (Watcher, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}
	if provider == nil {
		return nil, ErrNilProvider
	}
	if http == nil {
		return nil, ErrNilHTTP
	}
	domain := viper.GetStringMap("watcher")["domain"].(string)
	subdomains := utils.AggregateSubdomains(viper.GetStringMap("watcher")["subdomains"].([]string), domain)

	return &watcher{
		logger:     logger,
		provider:   provider,
		http:       http,
		domain:     domain,
		subdomains: subdomains,
	}, nil
}

func (w *watcher) Run() error {
	return nil
}
