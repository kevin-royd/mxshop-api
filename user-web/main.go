package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	// 1. 初始化zap日志
	initialize.InitLogger()

	// 2.获取配置文件
	initialize.InitConfig()

	// 3. 初始化svc客户端连接
	initialize.InitUserClient()

	// 4.初始化翻译器
	if err := initialize.InitValidator("zh"); err != nil {
		fmt.Printf("初始化翻译器错误, err = %s", err.Error())
		return
	}

	// 5.初始化router
	routers := initialize.Routers()
	if err := routers.Run(fmt.Sprintf(":%d", global.ServerConf.ServerPort)); err != nil {
		zap.S().Panicw("service start error", "msg", err.Error())
	}
}
