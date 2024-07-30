package main

import (
	ddnsclient "github.com/datahearth/ddnsclient"
	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = cobra.Command{
		Use:   "ddnsclient",
		Short: "ddnsclient is a dynamic DNS updater with built-in providers",
		Long: `ddnsclient will use a config file to update your A DNS settings periodicly.
						Checkout the documentation for parameters in the yaml config file.
					`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := ddnsclient.Start(logger, config); err != nil {
				logrus.Error(err)
			}
		},
	}
	logger = logrus.StandardLogger()
	config ddnsclient.ClientConfig
)

func init() {
	viper.BindEnv("CONFIG_PATH")
	viper.SetConfigType("yaml")
	if conf := viper.GetString("CONFIG_PATH"); conf == "" {
		viper.SetConfigFile("ddnsclient.yaml")
	} else {
		viper.SetConfigFile(conf)
	}

	utils.LoadConfig()
	if err := viper.Unmarshal(&config); err != nil {
		logger.WithError(err).Fatalln("failed to map yaml config file into ClientConfig struct")
	}

	utils.SetupLogger(logger)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.WithError(err).Fatalln("failed to execute command")
	}
}
