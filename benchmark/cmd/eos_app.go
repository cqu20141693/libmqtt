package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/server"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"net/http"
)

var Router = gin.Default()

func main() {
	api := Router.Group("/api/geega")
	MockPushData(api)
	CreateClientRoutes(api)
	GetClients(api)
	//server.StartServer(router)
	server.StartWithContextNotify(Router)
}
func MockPushData(rg *gin.RouterGroup) {

	rg.POST("/mockPush", func(c *gin.Context) {
		var topic = "v1/gateway/telemetry"

		var info map[string]interface{}
		err := c.BindJSON(&info)
		if err != nil {
			cclog.SugarLogger.Error(err)
			c.JSON(http.StatusBadRequest, "请求参数错误")
			return
		}
		cclog.SugarLogger.Info(fmt.Sprintf("topic: %s, data: %s", topic, info))

		for _, clientInfo := range domain.ClientMaps {
			client := clientInfo.GetClient()
			telemetryPktBytes, err := sonic.Marshal(info)
			if err != nil {
				cclog.SugarLogger.Errorf("sonic marchal faild: %v", err)
				return
			}
			client.Publish(&libmqtt.PublishPacket{
				TopicName: topic,
				Qos:       0,
				Payload:   telemetryPktBytes,
			})
		}
		c.JSON(http.StatusOK, "success")
	})
}
func CreateClientRoutes(rg *gin.RouterGroup) {

	rg.POST("/createClient", func(c *gin.Context) {
		var info domain.MqttClientAddInfo
		err := c.BindJSON(&info)
		if err != nil {
			cclog.SugarLogger.Error(err)
			c.JSON(http.StatusBadRequest, "请求参数错误")
			return
		}
		cclog.SugarLogger.Info(info)
		client, err := mqtt.CreatClient(&info)
		if err != nil {
			cclog.SugarLogger.Error(err)
			c.JSON(http.StatusInternalServerError, "libmqtt 连接失败")
			return
		}
		clientInfo := domain.NewGClientInfo(info.Address, info.ClientID, info.Username, info.Password, info.Keepalive)
		clientInfo.Connected = true
		clientInfo.SetClient(client)
		clientInfo.MockPolicy = info.MockPolicy
		domain.ClientMaps[clientInfo.ClientID] = clientInfo
		c.JSON(http.StatusOK, "success")
	})
}

func GetClients(rg *gin.RouterGroup) {
	rg.GET("/getClients", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			c.JSON(http.StatusOK, info)
		} else {
			c.JSON(http.StatusOK, domain.ClientMaps)
		}
	})
}
