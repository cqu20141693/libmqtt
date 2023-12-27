package gaea

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/domain"
	"github.com/goiiot/libmqtt/cmd/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/cmd/utils"
	"net/http"
	"time"
)

func StartMock(rg *gin.RouterGroup) {
	rg.GET("/StartMock", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			if info != nil {
				cclog.Info("start mock %v", info)
				startMock(info)
			}
		}
		c.JSON(http.StatusOK, "success")

	})
}
func StopMock(rg *gin.RouterGroup) {
	rg.GET("/StopMock", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			info.Scheduler().Stop()
			c.JSON(http.StatusOK, "success")
		} else {
			c.JSON(http.StatusOK, "success")
		}
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

func DisconnectClient(rg *gin.RouterGroup) {
	rg.GET("/disconnectClient", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			success := info.Client().Disconnect(info.Server, nil)
			c.JSON(http.StatusOK, fmt.Sprintf("disconnectClient %v", success))
		} else {
			c.JSON(http.StatusOK, "client not exist")
		}
	})
}

func ReconnectClient(rg *gin.RouterGroup) {
	rg.GET("/reconnectClient", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			err := info.Client().ReconnectServer(info.Server)
			if err != nil {
				cclog.Error(err)
				return
			}
			c.JSON(http.StatusOK, "success")
		} else {
			c.JSON(http.StatusOK, "client not exist")
		}
	})
}
func PublishMsg(rg *gin.RouterGroup) {
	rg.POST("/publishMsg", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			type PublishInfo struct {
				Topic   string                 `form:"topic" json:"topic" xml:"topic" binding:"required"`
				Qos     byte                   `form:"qos" json:"qos" xml:"qos" binding:"required"`
				Message map[string]interface{} `form:"message" json:"message" xml:"message" binding:"required"`
			}
			var infos []PublishInfo
			err := c.BindJSON(&infos)
			if err != nil {
				cclog.Error(err)
				c.JSON(http.StatusBadRequest, "请求参数错误")
				return
			}
			if info, ok := domain.ClientMaps[clientId]; ok {
				for i := range infos {
					msg, err := json.Marshal(infos[i].Message)
					if err != nil {
						cclog.Error(err)
						continue
					}
					message := string(msg)
					info.Client().Publish(utils.CreatePublishPacket(infos[i].Topic, infos[i].Qos, message))
				}
			}

			c.JSON(http.StatusOK, "success")
		} else {
			c.JSON(http.StatusOK, "client not exist")
		}
	})
}
func CreateClientRoutes(rg *gin.RouterGroup) {

	rg.POST("/createClient", func(c *gin.Context) {
		var info domain.MqttClientAddInfo
		err := c.BindJSON(&info)
		if err != nil {
			cclog.Error(err)
			c.JSON(http.StatusBadRequest, "请求参数错误")
			return
		}
		cclog.Info(info)
		options := make([]mqtt.Option, 0)
		options = append(options, mqtt.WithCleanSession(true))
		options = append(options, mqtt.WithClientID(info.ClientID))
		options = append(options, mqtt.WithIdentity(info.Username, info.Password))
		options = append(options, mqtt.WithKeepalive(uint16(info.Keepalive), 1.2))
		options = append(options, mqtt.WithVersion(mqtt.V311, false))
		client, err := newClient(options, info.Address)
		if err != nil {
			cclog.Error(err)
			c.JSON(http.StatusInternalServerError, "mqtt 连接失败")
			return
		}
		clientInfo := domain.NewGClientInfo(info.Address, info.ClientID, info.Username, info.Password, info.Keepalive)
		clientInfo.SetClient(client)
		clientInfo.MockPolicy = info.MockPolicy
		if startMock(clientInfo) {
			return
		}
		domain.ClientMaps[clientInfo.ClientID()] = clientInfo
		c.JSON(http.StatusOK, "success")
	})
}

func startMock(clientInfo *domain.GClientInfo) bool {
	if clientInfo.MockPolicy != nil {
		for i := range clientInfo.MockPolicy {
			if clientInfo.MockPolicy[i].Enable {
				policy := clientInfo.MockPolicy[i]
				scheduler := gocron.NewScheduler(time.Local)
				_, err := scheduler.Every(policy.Frequency).Millisecond().Do(func() {
					if clientInfo.EnableMock() {
						fmt.Printf("timer %v ms publish message", policy.Frequency)
						clientInfo.Client().Publish(utils.CreatePublishPacket(policy.Topic, policy.Qos, policy.Message))
					}

				})
				if err != nil {
					cclog.Error(err)
					return true
				}
				scheduler.StartAsync()
				clientInfo.SetScheduler(scheduler)
			}
		}
	}
	return false
}
