package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfig will read the yaml config from the viper config path
func LoadConfig() {
	logger := logrus.WithFields(logrus.Fields{
		"pkg":       "utils",
		"component": "config",
	})

	if err := viper.ReadInConfig(); err != nil {
		logger.WithError(err).Fatalln("failed to load configuration file")
	}
}

// SetupLogger setup the root logger
func SetupLogger(logger *logrus.Logger) {
	var (
		level        = logrus.InfoLevel
		timestamp    = true
		color        = true
		loggerConfig = viper.GetStringMap("logger")
	)

	if l, ok := loggerConfig["level"]; ok {
		parsedLevel, err := logrus.ParseLevel(l.(string))
		if err != nil {
			level = logrus.InfoLevel
		}
		level = parsedLevel
	}

	if t, ok := loggerConfig["disable-timestamp"]; ok {
		timestamp = t.(bool)
	}
	if c, ok := loggerConfig["disable-color"]; ok {
		color = c.(bool)
	}
	_ = timestamp

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    color,
		ForceColors:      true,
		FullTimestamp:    true,
		DisableTimestamp: timestamp,
	})
}

// RetrieveServerIP will use the defined web-ip service to get the server public address
func RetrieveServerIP(webIP string) (string, error) {
	resp, err := http.Get(webIP)
	if err != nil {
		return "", ErrGetServerIP
	}
	if resp.StatusCode != 200 {
		return "", ErrWrongStatusCode
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrParseHTTPBody
	}

	return string(d), nil
}

// RetrieveSubdomainIP will retrieve the subdomain IP
func RetrieveSubdomainIP(addr string) (string, error) {
	ips, err := net.LookupIP(addr)
	if err != nil {
		return "", err
	}

	if len(ips) != 1 {
		return "", ErrIpLenght
	}

	ip := ips[0].String()
	if strings.Contains(ip, ":") {
		ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			return "", ErrSplitAddr
		}
	}

	return ip, nil
}
