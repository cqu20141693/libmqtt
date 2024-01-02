package main

import (
	"github.com/goiiot/libmqtt/gateway/api"
	"github.com/goiiot/libmqtt/gateway/initialize/server"
)

func main() {
	api.InitLogger()
	// 启动api服务
	api.GetRoutes()
	//server.StartServer(router)
	server.StartWithContextNotify(api.Router)
	// todo 读取配置启动北向平台(mqtt server)

}
