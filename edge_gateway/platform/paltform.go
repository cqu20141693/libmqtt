package platform

import (
	"fmt"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/config"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/server"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
	"strings"
)

var Gateway = domain.Gateway{}

type Event struct {
	Ts         int64
	EventTopic events.EventTopicType // 事件主题
	DeviceId   string                // 事件目标设备
	Data       interface{}           //事件数据
}

type Telemetry Event

type DeviceEvent Event

type NorthPlatformEvent Event

type Reply Event

type SouthClient interface {

	// PublishEvent 推送事件
	//  @param event
	PublishEvent(event Event)

	// PublishTelemetry 推送数据
	//  @param event
	PublishTelemetry(event Telemetry)

	// PublishDeviceEvent 推送设备事件
	//  @param event
	PublishDeviceEvent(event DeviceEvent)

	// PublishReply 发送业务回复
	//  @param event
	PublishReply(event Reply)

	// ReceiveEvent 接收平台事件
	//  @return chan
	ReceiveEvent() chan NorthPlatformEvent
}

var PlatformClientMaps = make(map[string]SouthClient, 8)

func InitPlatform() {
	gatewayConfig := config.Viper.Get("gateway")
	if gatewayConfig != nil {
		cclog.SugarLogger.Info("read config success，", gatewayConfig)
		err := mapstructure.Decode(gatewayConfig, &Gateway)
		if err != nil {
			log.Fatal("decode edge_gateway config failed, ", err)
		}
		if Gateway.RestPort != 0 {
			server.DefaultAddr = ":" + strconv.Itoa(Gateway.RestPort)
		}
		if Gateway.AutoConnect {
			connectPlatform(Gateway)
		}

	}
}

func connectPlatform(g domain.Gateway) {
	platform := g.Platform

	switch platform.Type {
	case domain.GEEGA:
		address := strings.Join([]string{platform.Host, strconv.Itoa(platform.Port)}, ":")
		info := domain.NewMqttClientAddInfo(address, platform.ClientId, platform.Username, platform.Password, platform.KeepAlive)
		client, err := mqtt.CreatClient(info)
		if err != nil {
			cclog.Error(fmt.Sprintf("creat g platform mqtt client failed, %v", info), err.Error())

		}
		clientInfo := &GPlatformClient{
			Client:   client,
			ClientId: platform.ClientId,
		}
		PlatformClientMaps[clientInfo.ClientId] = clientInfo
	default:
		cclog.Warn(fmt.Sprintf("not suppport platform %v", platform))
	}

}
