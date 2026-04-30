package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "进入交互式命令行模式",
	Long:  `以 REPL 方式连续输入命令，直接与服务器交互。输入 help 查看可用命令，exit 退出。`,
	Run: func(cmd *cobra.Command, args []string) {
		startREPL()
	},
}

func startREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Axiom CLI 交互模式")
	fmt.Println("  输入 help 查看命令，exit 退出")
	fmt.Println("═══════════════════════════════════════")

	for {
		fmt.Print("axiom> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		switch args[0] {
		case "exit", "quit":
			fmt.Println("再见！")
			return

		case "help":
			printHelp()

		case "status":
			statusCmd.Run(&cobra.Command{}, args[1:])

		case "device":
			if len(args) > 1 && args[1] == "list" {
				deviceListCmd.Run(&cobra.Command{}, args[2:])
			} else {
				fmt.Println("用法: device list")
			}

		case "signal":
			if len(args) > 1 && args[1] == "health" {
				signalHealthCmd.Run(&cobra.Command{}, args[2:])
			} else {
				fmt.Println("用法: signal health")
			}

		case "rpc":
			if len(args) > 1 && args[1] == "list" {
				rpcListCmd.Run(&cobra.Command{}, args[2:])
			} else {
				fmt.Println("用法: rpc list")
			}

		case "clear":
			fmt.Print("\033[H\033[2J")

		default:
			fmt.Printf("未知命令: %s（输入 help 查看可用命令）\n", args[0])
		}

		if args[0] != "clear" {
			fmt.Println()
		}
	}
}

func printHelp() {
	fmt.Println("可用命令:")
	fmt.Println("  status             查看服务器状态")
	fmt.Println("  device list        列出设备")
	fmt.Println("  signal health      信令服务器健康检查")
	fmt.Println("  rpc list           列出 gRPC 接口")
	fmt.Println("  clear              清屏")
	fmt.Println("  exit / quit        退出")
	fmt.Println("  help               帮助")
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
