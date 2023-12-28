package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/cmd/initialize/server"
	"github.com/goiiot/libmqtt/cmd/schedule"
	"github.com/goiiot/libmqtt/cmd/service/gaea"
)

var router = gin.Default()

func main() {
	getRoutes()
	//server.StartServer(router)
	defer func() {
		schedule.Cron.Stop()
		schedule.GoCron.Stop()
	}()
	schedule.Cron.Start()
	// 异步启动
	schedule.GoCron.StartAsync()
	server.StartWithContextNotify(router)
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes() {

	//api
	initGMqttApi()
}

func initGMqttApi() {
	api := router.Group("/api/gmqtt")
	// 执行配置，建立client连接
	gaea.CreateClientRoutes(api)
	// 获取当前的客户端配置
	gaea.GetClients(api)
	//启用mock配置，开启定时推送任务
	gaea.StopMock(api)
	// 停止mock任务
	gaea.StartMock(api)
	// 更新mock任务
	gaea.UpdateMock(api)
	// 主动断开连接
	gaea.DisconnectClient(api)
	// 重连客户端
	gaea.ReconnectClient(api)
	// 一次性推送mqtt
	gaea.PublishMsg(api)

	// test api
	v1 := router.Group("/v1")
	gaea.AddPingRoutes(v1)

	v2 := router.Group("/v2")
	gaea.AddPingRoutes(v2)
}
