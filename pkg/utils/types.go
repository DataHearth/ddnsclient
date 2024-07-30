package utils

import (
	"errors"
)

// ** ERRORS ** //
var (
	// ErrNilLogger is thrown when the parameter logger is nil
	ErrNilLogger = errors.New("logger is mandatory")
	// ErrWrongStatusCode is thrown when the response status code isn't a 200
	ErrWrongStatusCode = errors.New("web-ip returns a non 200 status code")
	// ErrGetServerIP is thrown when HTTP can't contact the web-ip service
	ErrGetServerIP = errors.New("HTTP error")
	// ErrParseHTTPBody is thrown when the HTTP service can't parse the body response
	ErrParseHTTPBody = errors.New("can't parse response body")
	// ErrSplitAddr is thrown when the remote address can't be splitted
	ErrSplitAddr = errors.New("can't split subdomain remote IP address")
	// ErrCreateNewRequest is thrown when http request creation failed
	ErrCreateNewRequest = errors.New("can't create http request")
	// ErrInvalidURL is thrown when user does not provide a URL and it does not exist in default urls
	ErrInvalidURL = errors.New("no url was provided")
	// ErrInvalidName is thrown when provider name was not provided
	ErrInvalidName = errors.New("no provider name was provided")
	// ErrNilWatcher is thrown when no watcher config was provided
	ErrNilWatcher = errors.New("watcher is mandatory")
	// ErrIpLength is thrown when subdomain no or multiples remote IP address
	ErrIpLenght = errors.New("zero or more than 1 ips have been found")
	// ErrNilConfig is thrown when an empty config is provided
	ErrNilConfig = errors.New("config is mandatory")
)

// ** CONFIGURATION ** //
type ClientConfig struct {
	Logger     Logger    `mapstructure:"logger"`
	Watchers   []Watcher `mapstructure:"watchers"`
	UpdateTime int       `mapstructure:"update-time,omitempty"`
	WebIP      string    `mapstructure:"web-ip,omitempty"`
}

type Logger struct {
	Level            string `mapstructure:"level"`
	DisableTimestamp bool   `mapstructure:"disable-timestamp,omitempty"`
	DisableColor     bool   `mapstructure:"disable-color,omitempty"`
}

type Watcher struct {
	Name   string   `yaml:"name"`
	URL    string   `yaml:"url,omitempty"`
	Config []Config `yaml:"config"`
}

type Config struct {
	Username   string   `yaml:"username,omitempty"`
	Password   string   `yaml:"password,omitempty"`
	Token      string   `yaml:"password,omitempty"`
	Subdomains []string `yaml:"subdomains"`
}

var DefaultURLs = map[string]string{
	"ovh":     "http://www.ovh.com/nic/update?system=dyndns&hostname=SUBDOMAIN&myip=NEWIP",
	"google":  "https://domains.google.com/nic/update?hostname=SUBDOMAIN&myip=NEWIP",
	"duckdns": "https://duckdns.org/update/SUBDOMAIN/TOKEN[/NEWIP",
	"webIP":   "http://dynamicdns.park-your-domain.com/getip",
}
