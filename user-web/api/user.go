package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop-api/user-web/global"

	"mxshop-api/user-web/proto"
	"net/http"
)

var conn *grpc.ClientConn
var userClient proto.UserClient

func InitUserClient() {
	var err error
	conn, err = grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConf.UserServerInfo.Host, global.ServerConf.UserServerInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Error("【InitUserClient】连接用户服务失败：", err)
		return
	}
	userClient = proto.NewUserClient(conn)
}

func GetUserList(ctx *gin.Context) {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		global.HandleGrpcErrToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, rsp.Data)
}
