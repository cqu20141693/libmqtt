package main

import (
	"github.com/goiiot/libmqtt/edge_gateway/api"
	"github.com/goiiot/libmqtt/edge_gateway/gateway"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/server"
	"github.com/goiiot/libmqtt/edge_gateway/orm"
	"github.com/goiiot/libmqtt/edge_gateway/platform"
)

func main() {
	// 启动api服务
	api.GetRoutes()

	//  读取配置启动北向平台(mqtt server)
	platform.InitPlatform()

	// 初始化数据库
	orm.Init()

	//  初始化网关服务组件
	gateway.StartGateway()

	//server.StartServer(router)
	server.StartWithContextNotify(api.Router)
}
