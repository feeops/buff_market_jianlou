package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-version"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	// 导入grpc包 .
	"google.golang.org/grpc"
	// 导入刚才我们生成的代码所在的proto包 .
	pb "buff_market_jianlou/proto"
)

var (
	ProductName    = "buff_market_jianlou"
	ProductVersion = "1.0"
)

func auth(buffID string) error {
	// 连接grpc服务器
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		"gw.bufftools.com:443",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)

	if err != nil {
		return fmt.Errorf("无法连接服务器,错误原因:%w 请重试或切换网络", err)
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	// 初始化服务客户端
	c := pb.NewAccountClient(conn)

	// 发送消息
	req := new(pb.Request)
	req.ProductName = ProductName
	req.ProductVersion = ProductVersion
	req.BuffID = buffID

	r, err := c.Auth(ctx, req)
	if err != nil {
		return fmt.Errorf("验证失败::%w", err)
	}

	logger.Info().Str("ExpireTime", r.ExpireTime).Str("AccountID", r.AccountID).
		Str("currentProductVersion", ProductVersion).
		Str("latestProductVersion", r.ProductVersion).Msg("response")

	_ = os.WriteFile("账号信息.txt", []byte(r.AccountID), 0666)

	URL := fmt.Sprintf(`<meta http-equiv="refresh" content="0;url=https://www.bufftools.com/buy?accountID=%s">`,
		r.AccountID)
	_ = os.WriteFile("续费链接.html", []byte(URL), 0666)

	fmt.Printf("当前软件版本: %s 最新软件版本: %s\n", ProductVersion, r.ProductVersion)
	fmt.Printf("当前账号ID: %s 过期时间: %s\n", r.AccountID, r.ExpireTime)
	_ = os.Stdout.Sync()

	clientVersion, _ := version.NewVersion(ProductVersion)
	latestVersion, _ := version.NewVersion(r.ProductVersion)
	minVersion, _ := version.NewVersion(r.MinVersion)

	if clientVersion.LessThan(latestVersion) {
		fmt.Println("官网https://www.bufftools.com有最新的版本")
	}

	if len(r.MinVersion) > 0 && clientVersion.LessThan(minVersion) {
		alert := fmt.Sprintf("当前版本:%s 不满足最低版本:%s 请去官网https://www.bufftools.com下载最新版本\n",
			clientVersion, minVersion)
		fmt.Println(alert)
		waitExit()
	}

	if len(r.Notice) > 0 {
		alert := fmt.Sprintf("消息通知:%s 如有问题请联系作者\n",
			r.Notice)
		fmt.Println(alert)
		waitExit()
	}

	if r.ExpireSec <= 0 {
		alert := fmt.Sprintf("试用期：%s已过，请及时续费，续费后需要重启软件", r.ExpireTime)
		fmt.Println(alert)
		waitExit()
	}

	return nil
}
