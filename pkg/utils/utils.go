package utils

import (
	"io/ioutil"
	"net/http"

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

func AggregateSubdomains(subdomains []string, domain string) []string {
	agdSub := make([]string, len(subdomains))
	for i, sd := range subdomains {
		agdSub[i] = sd + "." + domain
	}

	return agdSub
}

// RetrieveServerIP will use the defined web-ip service to get the server public address and save it to the struct
func RetrieveServerIP(webIP string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"pkg":       "utils",
		"component": "server-ip",
	})

	// * retrieve client's server IP
	resp, err := http.Get(webIP)
	if err != nil {
		logger.WithError(err).WithField("web-ip", webIP).Errorln(ErrGetServerIP.Error())
		return "", ErrGetServerIP
	}
	if resp.StatusCode != 200 {
		logger.WithError(err).WithFields(logrus.Fields{
			"web-ip":      webIP,
			"statuc-code": resp.StatusCode,
		}).Errorln(ErrWrongStatusCode.Error())
		return "", ErrWrongStatusCode
	}

	// * get ip from body
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).WithField("web-ip", webIP).Errorln(ErrParseHTTPBody.Error())
		return "", ErrParseHTTPBody
	}

	return string(d), nil
}
