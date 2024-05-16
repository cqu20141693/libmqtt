package api

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"net/http"
	"time"
)

func StartMock(rg *gin.RouterGroup) {
	rg.GET("/StartMock", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			if info != nil {
				cclog.SugarLogger.Info("start mock %v", info)
				for _, policy := range info.MockPolicy {
					policy.Enable = true
				}
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
			err := info.Scheduler.Shutdown()
			if err != nil {
				cclog.SugarLogger.Errorf("Scheduler Shutdown failed=%v", err)
			}
			c.JSON(http.StatusOK, "success")
		} else {
			c.JSON(http.StatusOK, "success")
		}
	})
}
func UpdateMock(rg *gin.RouterGroup) {
	rg.POST("/updateMock", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			info := domain.ClientMaps[clientId]
			if info != nil {
				cclog.SugarLogger.Info("update mock %v", clientId)
				var policies []domain.PublishMockPolicy
				err := c.BindJSON(&policies)
				if err != nil {
					cclog.SugarLogger.Error(err)
					c.JSON(http.StatusBadRequest, "parameter error")
					return
				}
				info.MockPolicy = policies
				startMock(info)
			}
		}
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

func DisconnectClient(rg *gin.RouterGroup) {
	rg.GET("/disconnectClient", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			if info, ok := domain.ClientMaps[clientId]; ok {

				success := info.Client.Disconnect(info.Server, nil)
				info.Connected = false
				c.JSON(http.StatusOK, fmt.Sprintf("disconnectClient %v", success))
				return
			}
		}
		c.JSON(http.StatusOK, "client not exist")

	})
}

func ReconnectClient(rg *gin.RouterGroup) {
	rg.GET("/reconnectClient", func(c *gin.Context) {
		clientId, exist := c.GetQuery("clientId")
		if exist {
			if info, ok := domain.ClientMaps[clientId]; ok {
				if !info.Connected {
					err := info.Client.ReconnectServer(info.Server)
					if err != nil {
						cclog.SugarLogger.Error(err)
						c.JSON(http.StatusOK, fmt.Sprintf("reconnect failed %v", err))
						return
					}
					info.Connected = true
				} else {
					c.JSON(http.StatusOK, fmt.Sprintf("client already connected"))
					return
				}
				c.JSON(http.StatusOK, "success")
				return
			}
		}
		c.JSON(http.StatusOK, "client not exist")

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
				cclog.SugarLogger.Error(err)
				c.JSON(http.StatusBadRequest, "请求参数错误")
				return
			}
			if info, ok := domain.ClientMaps[clientId]; ok {
				for i := range infos {
					msg, err := json.Marshal(infos[i].Message)
					if err != nil {
						cclog.SugarLogger.Error(err)
						continue
					}
					message := string(msg)
					info.Client.Publish(utils.CreatePublishPacket(infos[i].Topic, infos[i].Qos, message))
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
		if startMock(clientInfo) {
			return
		}
		domain.ClientMaps[clientInfo.ClientID] = clientInfo
		c.JSON(http.StatusOK, "success")
	})
}

func startMock(clientInfo *domain.GClientInfo) bool {
	if clientInfo.EnableMock && clientInfo.MockPolicy != nil {
		for i := range clientInfo.MockPolicy {
			if clientInfo.MockPolicy[i].Enable {
				policy := clientInfo.MockPolicy[i]
				scheduler, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
				if err != nil {
					cclog.SugarLogger.Errorf("NewScheduler failed=%v", err)
					return true
				}
				_, err = scheduler.NewJob(
					gocron.DurationJob(time.Millisecond*time.Duration(policy.Frequency)),
					gocron.NewTask(func() {
						fmt.Printf("timer %v ms publish message", policy.Frequency)
						clientInfo.Client.Publish(utils.CreatePublishPacket(policy.Topic, policy.Qos, policy.Message))

					}),
				)
				if err != nil {
					cclog.SugarLogger.Error(err)
					return true
				}
				scheduler.Start()
				clientInfo.SetScheduler(scheduler)
			}
		}
	}
	return false
}
