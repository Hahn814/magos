/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logLevel = new(slog.LevelVar) // INFO by default
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

var cfgFile = ""

var RootCmd = &cobra.Command{
	Use:   "magosctl",
	Short: "Magos command line interface",
	Long:  ``, // TODO: write up a long description for cli
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.magos.yaml)")
	initEnvironment()
}

func initEnvironment() {
	logLevel.Set(slog.LevelDebug) // TODO: bind log level to environment and RootCmd

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".magos")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		logger.Debug("viper configuration read", "config_file", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Debug("config file not found")
		} else {
			logger.Error("an error occurred parsing the configuration file", "config_file", viper.ConfigFileUsed(), "error", err)
		}
	}

}
