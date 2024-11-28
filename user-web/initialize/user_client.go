package initialize

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitUserClient() {
	//conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConf.UserServerInfo.Host, global.ServerConf.UserServerInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 通过consul 注册resolver 添加lb算法
	// consul://[user:password@]127.0.0.127:8555/my-service?[healthy=]&[wait=]&[near=]&[insecure=]&[limit=]&[tag=]&[token=]

	cfg := global.ServerConf.ConsulInfo
	conn, err := grpc.NewClient(fmt.Sprintf("consul://%s:%d/%s?wait=%s&tag=%s", cfg.Host, cfg.Port,
		cfg.TargetServerName, "14s", cfg.Target),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Panicw("Init Client Conn Err", "error", err.Error())
		return
	}
	UserClient := proto.NewUserClient(conn)

	// NewClient使用的是grpc的懒加载机制不会主动发起连接。所以主动发起连接进行检查
	// 测试lb是否生效
	//for i := 0; i < 10; i++ {
	//	_, err = UserClient.GetUserById(context.Background(), &proto.IdRequest{Id: 1})
	//	if err != nil {
	//		zap.S().Errorw("Request Failed", "error", err)
	//	} else {
	//		zap.S().Infow("Response Received")
	//	}
	//	time.Sleep(1 * time.Second)
	//}

	_, err = UserClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: 1,
	})
	if err != nil {
		zap.S().Panicw("Init UserClient Err", "error", err.Error())
		return
	}
	global.UserClient = UserClient
}
