package cmd

import (
	"context"
	"fmt"
	"time"

	pb "axiom/server/pkg/grpc/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var deviceAddr string

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "设备管理",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有设备",
	Long:  `连接到业务服务器，列出当前测试账号下的所有绑定设备。`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := deviceAddr
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

		auth := pb.NewAuthClient(conn)
		auth.SendVerificationCode(ctx, &pb.SendCodeRequest{
			Phone: "+8618610320705", CountryCode: "CN",
		})
		loginResp, err := auth.LoginWithCode(ctx, &pb.LoginRequest{
			Phone: "+8618610320705", VerificationCode: "123456",
		})
		if err != nil {
			fmt.Printf("❌ 登录失败: %v\n", err)
			return
		}

		dev := pb.NewDeviceClient(conn)
		resp, err := dev.ListDevices(ctx, &pb.ListDevicesRequest{},
			grpc.PerRPCCredentials(&tokenCred{loginResp.Token}))
		if err != nil {
			fmt.Printf("❌ 获取设备列表失败: %v\n", err)
			return
		}

		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("  设备列表 (共 %d 个)\n", resp.TotalCount)
		fmt.Println("═══════════════════════════════════════")
		for _, d := range resp.Devices {
			fmt.Printf("  %-20s %-10s %s\n", d.Name, d.Status, d.Id)
			fmt.Printf("  %-20s %-10s 电量:%.0f%%\n", "", d.IpAddress, d.BatteryLevel)
			fmt.Println("  ─────────────────────────────")
		}
	},
}

type tokenCred struct{ token string }

func (t *tokenCred) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer " + t.token}, nil
}
func (t *tokenCred) RequireTransportSecurity() bool { return false }

func init() {
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.PersistentFlags().StringVar(&deviceAddr, "addr", "localhost:50053", "服务器地址")
}
