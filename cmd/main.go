package main

import (
	"log"

	ddnsclient "github.com/datahearth/ddns-client"
	"github.com/datahearth/ddns-client/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = cobra.Command{
		Use:   "ddns-client",
		Short: "ddns-client is a dynamic DNS updater with built-in providers",
		Long: `ddns-client will use a config file to update your A DNS settings periodicly.
						Checkout the documentation for parameters in the yaml config file.
					`,
		Run: func(cmd *cobra.Command, args []string) {
			ddnsclient.Start(logger)
		},
	}
	logger = logrus.StandardLogger()
)

func init() {
	viper.BindEnv("CONFIG_PATH")
	viper.SetConfigType("yaml")
	if conf := viper.GetString("CONFIG_PATH"); conf == "" {
		viper.SetConfigFile("ddns-client.yaml")
	} else {
		viper.SetConfigFile(conf)
	}

	if err := utils.LoadConfig(logger); err != nil {
		log.Fatalf("failed to load config file: %v\n", err.Error())
	}

	utils.SetupLogger(logger)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.WithError(err).Fatalln("failed to execute command")
	}
}
