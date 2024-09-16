package initialize

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitUserClient() {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConf.UserServerInfo.Host, global.ServerConf.UserServerInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicw("Init Client Conn Err", "error", err.Error())
		return
	}
	UserClient := proto.NewUserClient(conn)

	// NewClient使用的是grpc的懒加载机制不会主动发起连接。所以主动发起连接进行检查
	_, err = UserClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: 1,
	})
	if err != nil {
		zap.S().Panicw("Init UserClient Err", "error", err.Error())
		return
	}
	global.UserClient = UserClient
}
