package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var rpcAddr string

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "查看服务器 gRPC 接口",
}

var rpcListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有可用 gRPC 接口",
	Long:  `通过 gRPC reflection 协议查询服务器注册的所有服务和方法。`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := rpcAddr
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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ref := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
		stream, err := ref.ServerReflectionInfo(ctx)
		if err != nil {
			fmt.Printf("⚠️  reflection 不可用，显示已知接口:\n")
			printKnownServices()
			return
		}

		req := &grpc_reflection_v1alpha.ServerReflectionRequest{
			MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
				ListServices: "*",
			},
		}
		if err := stream.Send(req); err != nil {
			fmt.Printf("⚠️  reflection 请求失败，显示已知接口:\n")
			printKnownServices()
			return
		}

		resp, err := stream.Recv()
		if err != nil {
			fmt.Printf("⚠️  reflection 响应失败，显示已知接口:\n")
			printKnownServices()
			return
		}

		fmt.Println("═══════════════════════════════════════")
		fmt.Println("  服务器 gRPC 接口")
		fmt.Println("═══════════════════════════════════════")
		for _, svc := range resp.GetListServicesResponse().Service {
			fmt.Printf("  📡 %s\n", svc.Name)
		}
		fmt.Println("═══════════════════════════════════════")
	},
}

func printKnownServices() {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Axiom 服务器 gRPC 接口")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  📡 axiom.v1.Auth")
	fmt.Println("       SendVerificationCode  (Unary)")
	fmt.Println("       LoginWithCode         (Unary)")
	fmt.Println("       DeviceLogin           (Unary)")
	fmt.Println("       Logout                (Unary)")
	fmt.Println("       ChangePhone           (Unary)")
	fmt.Println("       ValidateToken         (Unary)")
	fmt.Println("  📡 axiom.v1.Device")
	fmt.Println("       ListDevices           (Unary)")
	fmt.Println("       GetDevice             (Unary)")
	fmt.Println("       BindDevice            (Unary)")
	fmt.Println("       AddDevice             (Unary)")
	fmt.Println("       RemoveDevice          (Unary)")
	fmt.Println("       UpdateDevice          (Unary)")
	fmt.Println("       SubscribeDeviceStatus (ServerStream)")
	fmt.Println("       ReportDeviceOnline    (Unary)")
	fmt.Println("       ReportDeviceOffline   (Unary)")
	fmt.Println("       ReportDeviceStatus    (Unary)")
	fmt.Println("       GenerateDeviceBatch   (Unary)")
	fmt.Println("  📡 axiom.v1.RobotControl")
	fmt.Println("       CommandStream         (BidiStream)")
	fmt.Println("  📡 grpc.health.v1.Health")
	fmt.Println("       Check                 (Unary)")
	fmt.Println("       Watch                 (ServerStream)")
	fmt.Println("═══════════════════════════════════════")
}

func init() {
	rpcCmd.AddCommand(rpcListCmd)
	rpcCmd.PersistentFlags().StringVar(&rpcAddr, "addr", "localhost:50053", "服务器地址")
}
