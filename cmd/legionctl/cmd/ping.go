/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	legionpb "github.com/Hahn814/legion/proto/legion/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("legion")
		viper.BindEnv("port")
		viper.SetDefault("port", 50051)
		port := viper.GetInt("port")
		addr := viper.GetString("addr")
		if addr == "" {
			addr = "0.0.0.0"
		}

		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error("did not connect", "error", err)
			os.Exit(1)
		}
		defer conn.Close()
		c := legionpb.NewLegionClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.Hello(ctx, &legionpb.HelloRequest{Name: "legionctl"})
		if err != nil {
			logger.Error("could not ping %v", "error", err)
			os.Exit(1)
		}

		logger.Info(r.GetName())
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
