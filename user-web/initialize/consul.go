package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
)

// InitConsul 从consul获取service user-srv
func InitConsul() {
	config := api.DefaultConfig()
	ConsulInfo := global.ServerConf.ConsulInfo
	config.Address = fmt.Sprintf("%s:%d", ConsulInfo.Host, ConsulInfo.Port)
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panicw("init consul fail", "err", err)
	}
	// 通过agent filter获取service
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConf.ConsulInfo.TargetServerName))
	if err != nil {
		zap.S().Panicw("Filter consul Service fail", "err", err)
	}

	for _, value := range data {
		// 通过consul 获取grpc的地址
		global.ServerConf.UserServerInfo.Host = value.Address
		global.ServerConf.UserServerInfo.Port = value.Port
	}
	fmt.Printf("user-srv-info:%v\n", global.ServerConf.UserServerInfo)
}
