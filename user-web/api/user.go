package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/proto"
	"net/http"
)

// 将grpc状态码转换为http
func HandlerGrpcToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错位",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错位",
				})
			}
		}
	}
}

func GetUserList(ctx *gin.Context) {
	// 拨号连接用户grpc服务端
	userConn, err := grpc.NewClient("127.0.0.1:8088", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接服务端", "msg", err.Error())

	}
	// 调用接口
	userClient := proto.NewUserClient(userConn)
	rsp, err := userClient.GetUserList(ctx, &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询用户列表失败", "msg", err.Error())
		HandlerGrpcToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		data := make(map[string]interface{})
		data["id"] = value.Id
		data["mobile"] = value.Mobile
		data["nickname"] = value.Nickname
		data["gender"] = value.Gender
		data["birthDay"] = value.BirthDay
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
}
