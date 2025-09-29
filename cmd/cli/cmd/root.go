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

var RootCmd = &cobra.Command{
	Use:   "magosctl",
	Short: "Magos command line interface",
	Long:  `magosctl is the command line interface used to manage the Magos agents`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("verbosity", "v", "DEBUG", "logging verbosity (DEBUG, INFO, WARN, ERROR)")
	viper.BindPFlag("verbosity", RootCmd.PersistentFlags().Lookup("verbosity"))
	initEnvironment()
}

func initEnvironment() {
	viper.SetEnvPrefix("magos")
	viper.AutomaticEnv()

	// TODO: add this logging logic to a shared internal package to improve reusability
	levelName := viper.GetString("verbosity")
	err := logLevel.UnmarshalText([]byte(levelName))
	if err != nil {
		logger.Warn("failed to parse log level", "levelName", levelName, "error", err)
		logger.Info("defaulting log level to INFO")
		logLevel.Set(slog.LevelInfo)
	}

	viper.Set("verbosity", logLevel.Level().String())
	logger.Debug("log level configured", "level", viper.GetString("verbosity"))
}
