package ddnsclient

import (
	"github.com/datahearth/ddnsclient/pkg/providers/google"
	"github.com/datahearth/ddnsclient/pkg/providers/ovh"
)

type ClientConfig struct {
	Logger                Logger        `mapstructure:"logger"`
	Providers             Providers     `mapstructure:"providers"`
	Watchers              WatcherConfig `mapstructure:"watchers"`
	UpdateTime            int           `mapstructure:"update-time,omitempty"`
	PendingDnsPropagation int           `mapstructure:"pending-dns-propagation,omitempty"`
	WebIP                 string        `mapstructure:"web-ip,omitempty"`
}

type Logger struct {
	Level            string `mapstructure:"level"`
	DisableTimestamp bool   `mapstructure:"disable-timestamp,omitempty"`
	DisableColor     bool   `mapstructure:"disable-color,omitempty"`
}

type Providers struct {
	Ovh    ovh.OvhConfig       `mapstructure:"ovh,omitempty"`
	Google google.GoogleConfig `mapstructure:"google,omitempty"`
}

type WatcherConfig struct {
	Ovh    []string `yaml:"ovh,omitempty"`
	Google []string `yaml:"google,omitempty"`
}
