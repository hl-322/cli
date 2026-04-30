package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var signalAddr string

var signalCmd = &cobra.Command{
	Use:   "signal",
	Short: "信令服务器",
}

var signalHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "信令服务器健康检查",
	Run: func(cmd *cobra.Command, args []string) {
		addr := signalAddr
		if addr == "" {
			addr = "localhost:50052"
		}

		conn, err := grpc.NewClient(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())
		if err != nil {
			fmt.Printf("❌ 连接失败: %v\n", err)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := grpc_health_v1.NewHealthClient(conn)
		resp, err := health.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			fmt.Printf("❌ 信令服务器未响应: %v\n", err)
			return
		}

		fmt.Println("═══════════════════════════════════════")
		fmt.Println("  信令服务器")
		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("  状态:   %s\n", resp.Status.String())
		fmt.Printf("  地址:   %s\n", addr)
		fmt.Println("═══════════════════════════════════════")
	},
}

func init() {
	signalCmd.AddCommand(signalHealthCmd)
	signalCmd.PersistentFlags().StringVar(&signalAddr, "addr", "localhost:50052", "信令服务器地址")
}
