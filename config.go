package ddnsclient

type ClientConfig struct {
	Logger     Logger    `mapstructure:"logger"`
	Providers  Providers `mapstructure:"providers"`
	Watcher    Watcher   `mapstructure:"watcher"`
	UpdateTime int       `mapstructure:"update-time"`
	WebIP      string    `mapstructure:"web-ip"`
}

type Logger struct {
	Level            string `mapstructure:"level"`
	DisableTimestamp bool   `mapstructure:"disable-timestamp"`
	DisableColor     bool   `mapstructure:"disable-color"`
}

type Providers struct {
	Ovh Ovh `mapstructure:"ovh,omitempty"`
}

type Ovh struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Watcher struct {
	Domain     string   `mapstructure:"domain"`
	Subdomains []string `mapstructure:"subdomains"`
}
