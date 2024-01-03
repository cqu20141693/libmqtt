package platform

import (
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/gateway/initialize/config"
	"github.com/goiiot/libmqtt/gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/gateway/initialize/server"
	"github.com/goiiot/libmqtt/gateway/mqtt"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
	"strings"
)

var Gateway = domain.Gateway{}

func InitPlatformConfig() {
	gatewayConfig := config.Viper.Get("gateway")
	if gatewayConfig != nil {
		cclog.SugarLogger.Info("read config success，", gatewayConfig)
		config := gatewayConfig.(map[string]interface{})
		err := mapstructure.Decode(config, &Gateway)
		if err != nil {
			log.Fatal("decode gateway config failed, ", err)
		}
		if Gateway.RestPort != 0 {
			server.DefaultAddr = ":" + strconv.Itoa(Gateway.RestPort)
		}
		if Gateway.AutoConnect {
			if connectPlatformAndStartTask(Gateway) {
				return
			}
		}
	}
}

func connectPlatformAndStartTask(g domain.Gateway) bool {
	platform := g.Platform
	address := strings.Join([]string{platform.Host, strconv.Itoa(platform.Port)}, ":")
	info := domain.NewMqttClientAddInfo(address, platform.ClientId, platform.Username, platform.Password, platform.KeepAlive)
	client, err := mqtt.CreatClient(info)
	if err != nil {
		cclog.Error("creat g platform mqtt client failed, ", err.Error())
		return true
	}
	clientInfo := domain.NewGClientInfo(info.Address, info.ClientID, info.Username, info.Password, info.Keepalive)
	clientInfo.Connected = true
	clientInfo.SetClient(client)
	domain.PlatformClientMaps[clientInfo.ClientID] = clientInfo
	return false
}
