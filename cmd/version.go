package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.hollow.sh/toolbox/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
