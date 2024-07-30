package utils

type ClientConfig struct {
	Logger                Logger     `mapstructure:"logger"`
	Watchers              []Watcher `mapstructure:"watchers"`
	UpdateTime            int        `mapstructure:"update-time,omitempty"`
	PendingDnsPropagation int        `mapstructure:"pending-dns-propagation,omitempty"`
	WebIP                 string     `mapstructure:"web-ip,omitempty"`
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
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Subdomains []string `yaml:"subdomains"`
}

var DefaultURLs = map[string]string{
	"ovh":    "http://www.ovh.com/nic/update?system=dyndns&hostname=SUBDOMAIN&myip=NEWIP",
	"google": "https://domains.google.com/nic/update?hostname=SUBDOMAIN&myip=NEWIP",
	"webIP":  "http://dynamicdns.park-your-domain.com/getip",
}
