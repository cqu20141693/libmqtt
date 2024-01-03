package main

import (
	"github.com/goiiot/libmqtt/gateway/api"
	"github.com/goiiot/libmqtt/gateway/initialize/server"
	"github.com/goiiot/libmqtt/gateway/platform"
)

func main() {
	// 启动api服务
	api.GetRoutes()

	// todo 读取配置启动北向平台(mqtt server)
	// todo 初始化网关服务组件：事件channel,本地存储队列(当北向连接不上，数据具有一定当存储能力)
	// todo 读取配置启动南向连接设备(device connector)
	platform.InitPlatformConfig()

	//server.StartServer(router)
	server.StartWithContextNotify(api.Router)
}
