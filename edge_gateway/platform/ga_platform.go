package platform

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
)

const (
	onlineTopic    = "v1/gateway/online"
	offlineTopic   = "v1/gateway/offline"
	telemetryTopic = "v1/gateway/telemetry"
	attributeTopic = "v1/gateway/attributes"
	connectTopic   = "v1/gateway/connect"
)

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
	if event.Ts == 0 {
		tsData["ts"] = event.Ts
	} else {
		tsData["ts"] = utils.GetTimestamp()
	}
	tsDatas = append(tsDatas, tsData)
	info[event.DeviceId] = tsDatas
	marshal, _ := sonic.Marshal(info)
	switch event.EventTopic {
	case events.TelemetryTopic:
		G.Client.Publish(utils.CreatePublishPacket(telemetryTopic, 1, string(marshal)))
	case events.AttributeTopic:
		G.Client.Publish(utils.CreatePublishPacket(attributeTopic, 1, string(marshal)))
	}

}

func (G *GPlatformClient) PublishDeviceEvent(event DeviceEvent) {
	switch event.EventTopic {
	case events.OnlineTopic:
		info := map[string][]string{}
		deviceIds := event.Data.([]string)
		info["deviceIds"] = deviceIds
		marshal, _ := json.Marshal(info)
		G.Client.Publish(utils.CreatePublishPacket(onlineTopic, 1, string(marshal)))
	case events.OfflineTopic:
		info := map[string][]string{}
		deviceIds := event.Data.([]string)
		info["deviceIds"] = deviceIds
		marshal, _ := json.Marshal(info)
		G.Client.Publish(utils.CreatePublishPacket(offlineTopic, 1, string(marshal)))
	case events.RegisterTopic:
		devices := event.Data.([]connectors.Device)
		for _, device := range devices {
			info := getConnectInfo(device)
			marshal, _ := json.Marshal(info)
			G.Client.Publish(utils.CreatePublishPacket(connectTopic, 1, string(marshal)))
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

func (G *GPlatformClient) ReceiveEvent() chan NorthPlatformEvent {
	//TODO implement me
	panic("implement me")
}
