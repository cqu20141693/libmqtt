package platform

import (
	"encoding/json"
	"fmt"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/config"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/server"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
	"strings"
)

var Gateway = domain.Gateway{}

type Event struct {
	EventTopic events.EventTopicType // 事件主题
	DeviceId   string                // 事件目标设备
	Data       interface{}           //事件数据
}

type Telemetry Event

type DeviceEvent Event

type PlatformEvent Event

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
	ReceiveEvent() chan PlatformEvent
}
type GPlatformClient struct {
	Client   libmqtt.Client
	ClientId string
}

func (G *GPlatformClient) PublishReply(event Reply) {
	//TODO implement me
	panic("implement me")
}

func (G *GPlatformClient) PublishEvent(event Event) {
	//TODO implement me
	panic("implement me")
}

func (G *GPlatformClient) PublishTelemetry(event Telemetry) {
	var info = map[string]interface{}{}
	var tsDatas = make([]interface{}, 0)
	var tsData = make(map[string]interface{})

	data := event.Data.(map[string]interface{})
	tsData["values"] = data
	tsData["ts"] = utils.GetTimestamp()
	tsDatas = append(tsDatas, tsData)
	info[event.DeviceId] = tsDatas
	marshal, _ := json.Marshal(info)
	G.Client.Publish(utils.CreatePublishPacket("v1/gateway/telemetry", 1, string(marshal)))
}

func (G *GPlatformClient) PublishDeviceEvent(event DeviceEvent) {
	switch event.EventTopic {
	case events.OnlineTopic:
		info := map[string][]string{}
		deviceIds := event.Data.([]string)
		info["deviceIds"] = deviceIds
		marshal, _ := json.Marshal(info)
		G.Client.Publish(utils.CreatePublishPacket("v1/gateway/online", 1, string(marshal)))
	case events.OfflineTopic:
		info := map[string][]string{}
		deviceIds := event.Data.([]string)
		info["deviceIds"] = deviceIds
		marshal, _ := json.Marshal(info)
		G.Client.Publish(utils.CreatePublishPacket("v1/gateway/offline", 1, string(marshal)))
	case events.RegisterTopic:
		devices := event.Data.([]connectors.Device)
		for _, device := range devices {
			info := getConnectInfo(device)
			marshal, _ := json.Marshal(info)
			G.Client.Publish(utils.CreatePublishPacket("v1/gateway/connect", 1, string(marshal)))
		}

	}
}

func getConnectInfo(device connectors.Device) []interface{} {
	var infos = make([]interface{}, 0)

	var info = map[string]interface{}{}
	info["deviceId"] = device.DeviceID
	info["deviceName"] = device.DeviceName
	info["productId"] = device.DeviceType
	info["productName"] = device.DeviceTypeName
	var attrs = make([]map[string]string, 0)
	var telemetries = make([]map[string]string, 0)
	for _, tag := range device.Tags {
		var tagInfo = map[string]string{}
		tagInfo["name"] = tag.Name
		tagInfo["key"] = tag.TagId
		tagInfo["dataType"] = tag.DataType
		if tag.TagType == constants.TIMESERIES {
			telemetries = append(telemetries, tagInfo)
		} else {
			attrs = append(attrs, tagInfo)
		}
	}
	//info[constants.ATTRIBUTES] = attrs
	//info[constants.TIMESERIES] = telemetries
	infos = append(infos, info)

	return infos
}

func (G *GPlatformClient) ReceiveEvent() chan PlatformEvent {
	//TODO implement me
	panic("implement me")
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
