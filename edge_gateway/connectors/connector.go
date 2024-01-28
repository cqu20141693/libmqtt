package connectors

import (
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"github.com/mitchellh/mapstructure"
)

type ConnectorManager struct {
	NorthEventCh chan events.NorthEvent
	SouthEventCh chan events.SouthEvent
	// 通道缓存
	// channelId -> Channel
	ChannelInstances map[string]Channel

	//  Devices 设备缓存
	// deviceId-> channelId
	Devices map[string]string
}

func (m *ConnectorManager) HandleChannelAdd(events.NorthChannelEvent) {

}
func (m *ConnectorManager) HandleChannelStop() {

}
func (m *ConnectorManager) HandleChannelRestart() {

}

// start  启动服务
//
//	@receiver m
func (m *ConnectorManager) Start() {
	// 处理北向任务
	utils.GoWithRecover(m.HandleNorthEvent)

}

func (m *ConnectorManager) HandleNorthEvent() {
	select {
	case event, more := <-m.NorthEventCh:
		if !more {
			return
		}
		switch event.Type {
		case constants.FUNCTION:
			northEvent := event.Data.(events.NorthCmdEvent)
			if channel, ok := m.ChannelInstances[northEvent.ChannelId]; ok {
				var rpc RpcEvent
				err := mapstructure.Decode(northEvent.Command, &rpc)
				if err != nil {
					cclog.Warn("decode UpdateTagEvent failed:", northEvent.Command)
					break
				}
				channel.ServerSideRpcHandler(rpc)
			}
		case constants.ATTRIBUTES_READ:
			northEvent := event.Data.(events.NorthCmdEvent)
			if channel, ok := m.ChannelInstances[northEvent.ChannelId]; ok {
				var update UpdateTagEvent
				err := mapstructure.Decode(northEvent.Command, &update)
				if err != nil {
					cclog.Warn("decode UpdateTagEvent failed:", northEvent.Command)
					break
				}
				channel.AttributesUpdate(update)
			}
		case constants.ATTRIBUTES_UPDATE:
			northEvent := event.Data.(events.NorthCmdEvent)
			if channel, ok := m.ChannelInstances[northEvent.ChannelId]; ok {
				var readEvent ReadTagEvent
				err := mapstructure.Decode(northEvent.Command, &readEvent)
				if err != nil {
					cclog.Warn("decode ReadTagEvent failed:", northEvent.Command)
					break
				}
				channel.AttributesRead(readEvent)
			}
		case constants.ADD_CHANNEL:
			channelEvent := event.Data.(events.NorthChannelEvent)
			m.HandleChannelAdd(channelEvent)
		case constants.STOP_CHANNEL:

		}
	}
}

// DeviceTelemetry
//
//	@receiver m
//	@param ret
func (m *ConnectorManager) DeviceTelemetry(ret ReadTagResult) {
	event := events.SouthEvent{
		Type:       events.Telemetry,
		EventTopic: events.TelemetryTopic,
		DeviceId:   ret.DeviceId,
		Data:       ret.Success,
	}
	m.SouthEventCh <- event
}

// DeviceAttribute
//
//	@receiver m
//	@param ret
func (m *ConnectorManager) DeviceAttribute(ret ReadTagResult) {
	event := events.SouthEvent{
		Type:       events.Attribute,
		EventTopic: events.AttributeTopic,
		DeviceId:   ret.DeviceId,
		Data:       ret,
	}
	m.SouthEventCh <- event
}

// DeviceOnline 设备上线
//
//	@receiver m
//	@param devices 设备ID
func (m *ConnectorManager) DeviceOnline(devices []string) {

	event := events.SouthEvent{
		Type:       events.DeviceEvent,
		EventTopic: events.OnlineTopic,

		Data: devices,
	}
	m.SouthEventCh <- event
}

// RegisterDevices 注册设备，自动创建设备点位
//
//	@receiver m
//	@param devices 设备列表
func (m *ConnectorManager) RegisterDevices(devices []Device) {
	// 缓存设备
	for _, device := range devices {
		m.Devices[device.DeviceID] = device.ChannelID
	}

	// 发送设备事件
	event := events.SouthEvent{
		Type:       events.DeviceEvent,
		EventTopic: events.RegisterTopic,
		Data:       devices,
	}
	m.SouthEventCh <- event
}

func (m *ConnectorManager) DeviceOffline(deviceIds []string) {
	// 清理设备
	for _, device := range deviceIds {
		delete(m.Devices, device)
	}
	event := events.SouthEvent{
		Type:       events.DeviceEvent,
		EventTopic: events.OfflineTopic,

		Data: deviceIds,
	}
	m.SouthEventCh <- event
}

// AddChannel 根据配置创建通道并启动通道
func AddChannel(event domain.AddChannelEvent) {

}
