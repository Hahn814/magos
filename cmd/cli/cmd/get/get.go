package get

import (
	"fmt"

	"github.com/Hahn814/magos/cmd/cli/cmd"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(GetCmd)
}
