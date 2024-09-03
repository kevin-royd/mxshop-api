package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/initialize"
)

func main() {
	initialize.InitLogger()

	// 2.初始化路由
	port := 8081
	zap.S().Debugf("启动服务 端口 %d", port)
	routers := initialize.Routers()
	if err := routers.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
