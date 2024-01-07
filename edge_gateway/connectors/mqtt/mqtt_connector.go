package mqtt

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/mitchellh/mapstructure"
	"time"
)

type MqttConnector struct {
	connectors.ChannelInstance
	Config        map[string]interface{} // 整体配置
	ConnectConfig map[string]interface{} // 连接配置，根据协议不同，连接配置不同
}

func (m *MqttConnector) ServerSideRpcHandler(eventData connectors.RpcEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) ServerSideRpcHandlerSync(eventData connectors.RpcEvent) connectors.RpcResult {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) AttributesUpdate(eventData connectors.UpdateTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) AttributesUpdateSync(eventData connectors.UpdateTagEvent) connectors.UpdateTagResult {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) AttributesRead(eventData connectors.ReadTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) AttributesReadSync(eventData connectors.ReadTagEvent) connectors.ReadTagResult {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) Close(remote bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) IsConnected() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) GetId() string {
	//TODO implement me
	panic("implement me")
}

func (m *MqttConnector) Open() {

}
func (m *MqttConnector) Init() {
	// 初始化连接配置
	m.ConnectConfig = m.Config["connector"].(map[string]interface{})
	// 解析配置为ChannelCache
}

func NewMqttConnector(connectorManager *connectors.ConnectorManager, channelId string, channelName string, config map[string]interface{}) connectors.Channel {
	var reportConfig connectors.ReportConfig
	if conf, ok := config[constants.ReportConfigKey]; ok {
		err := mapstructure.Decode(conf, &reportConfig)
		if err != nil {
			cclog.Error(fmt.Sprintf("decode report config failed,%v %v", channelName, conf), err)
		}
	}
	//初始化Scheduler，ChannelCache，Connected
	var channelInstance = connectors.ChannelInstance{ConnectorManager: connectorManager, ChannelId: channelId,
		ChannelName: channelName, Scheduler: gocron.NewScheduler(time.Local),
		ChannelCache: connectors.ChannelCache{}, ReportConfig: reportConfig, Connected: false, Enabled: true}
	mqtt := &MqttConnector{ChannelInstance: channelInstance, Config: config}
	// 初始化配置
	mqtt.Init()
	return mqtt
}
