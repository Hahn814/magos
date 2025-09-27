package get

import (
	"log/slog"
	"os"

	"github.com/Hahn814/magos/cmd/cli/cmd/get"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		logger.Debug("agent subcommand called", "args", args)
		for _, agentId := range args {
			// TODO: use the API server to get the registered agent details
			logger.Debug("inspect", "agentId", agentId)
		}
	},
}

func init() {
	logLevel.UnmarshalText([]byte(viper.GetString("verbosity")))
	get.GetCmd.AddCommand(agentCmd)
}
