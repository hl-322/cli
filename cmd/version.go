package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	AppVersion = "1.0.0"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("axiom-cli v%s\n", AppVersion)
		fmt.Printf("Build: %s | Commit: %s\n", BuildTime, CommitHash)
	},
}
