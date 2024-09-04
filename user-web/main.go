package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	// 1. 初始化zap日志
	initialize.InitLogger()

	// 2.获取配置文件
	initialize.InitConfig()

	// 3. 初始化svc客户端连接
	api.InitUserClient()

	// 2.初始化router
	routers := initialize.Routers()
	if err := routers.Run(fmt.Sprintf(":%d", global.ServerConf.ServerPort)); err != nil {
		zap.S().Panicw("service start error", "msg", err.Error())
	}
}
