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
	api := router.Group("/api/gmqtt")
	gaea.CreateClientRoutes(api)
	gaea.GetClients(api)
	gaea.StopMock(api)
	gaea.StartMock(api)
	gaea.DisconnectClient(api)
	gaea.ReconnectClient(api)
	gaea.PublishMsg(api)
	v1 := router.Group("/v1")
	gaea.AddPingRoutes(v1)

	v2 := router.Group("/v2")
	gaea.AddPingRoutes(v2)
}
