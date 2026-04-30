package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "axiom",
	Short: "Axiom 机器人控制系统 CLI",
	Long:  `通过 gRPC 连接 Axiom 服务器进行管理和监控。`,
	Run: func(cmd *cobra.Command, args []string) {
		startREPL()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(signalCmd)
	rootCmd.AddCommand(rpcCmd)
	rootCmd.AddCommand(versionCmd)
}
