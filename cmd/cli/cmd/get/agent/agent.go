package get

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Hahn814/magos/cmd/cli/cmd/get"
	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logLevel = new(slog.LevelVar)
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent [agent_id]",
	Short: "List a single agents details",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.MinimumNArgs(1)(cmd, args)
		if err != nil {
			logger.Error("agent id argument required")
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: refactor acquiring api server to shared internal package
		logger.Debug("agent subcommand called", "args", args)
		for _, agentId := range args {
			logger.Debug("inspect", "agentId", agentId)
			dialOpts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}
			conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "0.0.0.0", 50051), dialOpts...)
			if err != nil {
				logger.Error("did not connect", "error", err)
				os.Exit(1)
			}
			defer conn.Close()

			client := magosapipb.NewAPIClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			r, err := client.GetAgent(ctx, &magostypespb.GetAgentRequest{Id: agentId})
			if err != nil {
				logger.Error("could not get agent", "error", err)
				os.Exit(2)
			}

			logger.Debug("get agent", "hostname", r.GetHostname(), "id", r.GetId()) // TODO: implement shell output adapter
		}
	},
}

func init() {
	logLevel.UnmarshalText([]byte(viper.GetString("verbosity")))
	get.GetCmd.AddCommand(agentCmd)
}
