package cmd

import (
	"context"
	"fmt"
	"time"

	pb "axiom/server/pkg/grpc/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var statusAddr string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看服务器状态",
	Long:  `连接到运行中的业务服务器，查看 gRPC 服务状态。`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := statusAddr
		if addr == "" {
			addr = "localhost:50053"
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
		fmt.Println("═══════════════════════════════════════")
		fmt.Println("  Axiom 业务服务器状态")
		fmt.Println("═══════════════════════════════════════")

		for _, svc := range []struct{ name, method string }{
			{"整体服务", ""},
			{"Auth", "axiom.v1.Auth"},
			{"Device", "axiom.v1.Device"},
			{"RobotControl", "axiom.v1.RobotControl"},
		} {
			resp, err := health.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: svc.method})
			if err != nil {
				fmt.Printf("  %-20s ❌ %v\n", svc.name, err)
			} else if resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
				fmt.Printf("  %-20s ✅ SERVING\n", svc.name)
			} else {
				fmt.Printf("  %-20s ❌ %s\n", svc.name, resp.Status)
			}
		}

		// 额外验证 Auth 服务
		auth := pb.NewAuthClient(conn)
		_, err = auth.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: "test"})
		if err != nil {
			fmt.Printf("\n  Auth 服务: ✅ 运行中\n")
		}

		fmt.Printf("\n  地址: %s\n", addr)
		fmt.Println("═══════════════════════════════════════")
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusAddr, "addr", "localhost:50053", "服务器地址")
}
