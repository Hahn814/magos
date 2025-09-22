/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"time"

	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var username = ""

// helloCmd represents the ping command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Submit a HelloRequest to the Magos agent service",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("magos")
		viper.BindEnv("api.port")
		viper.SetDefault("api.port", 50051)
		port := viper.GetInt("api.port")
		addr := viper.GetString("api.addr")
		if addr == "" {
			addr = "0.0.0.0"
		}

		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error("did not connect", "error", err)
			os.Exit(1)
		}
		defer conn.Close()
		c := magosapipb.NewAPIClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if username == "" {
			currentUser, err := user.Current()
			if err != nil {
				logger.Error("Error getting current user", "error", err)
				os.Exit(1)
			}

			username = currentUser.Username
		}

		r, err := c.Hello(ctx, &magostypespb.HelloRequest{Name: username})
		if err != nil {
			logger.Error("could not ping %v", "error", err)
			os.Exit(1)
		}

		logger.Info(r.GetName())
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	helloCmd.Flags().StringVar(&username, "name", "", "Alternate name to provide in the HelloRequest (Default is current OS user")
}
