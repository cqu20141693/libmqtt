package opcua

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/mitchellh/mapstructure"
	"time"
)

type OpcUaConnector struct {
	connectors.ChannelInstance
	Config        map[string]interface{} // 整体配置
	ConnectConfig map[string]interface{} // 连接配置，根据协议不同，连接配置不同
}

func (m *OpcUaConnector) ServerSideRpcHandler(eventData connectors.RpcEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) AttributesUpdate(eventData connectors.UpdateTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) AttributesRead(eventData connectors.ReadTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) ServerSideRpcHandlerSync(eventData connectors.RpcEvent) connectors.RpcResult {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) AttributesUpdateSync(eventData connectors.UpdateTagEvent) connectors.UpdateTagResult {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) AttributesReadSync(eventData connectors.ReadTagEvent) connectors.ReadTagResult {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) Open() {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) Close(remote bool) {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) IsConnected() bool {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) GetId() string {
	//TODO implement me
	panic("implement me")
}

func (m *OpcUaConnector) Init() {
	// 初始化连接配置
	m.ConnectConfig = m.Config["connector"].(map[string]interface{})
	// 解析配置为ChannelCache
}
func NewOpcUaConnector(connectorManager *connectors.ConnectorManager, channelId string, channelName string, config map[string]interface{}) connectors.Channel {
	var reportConfig connectors.ReportConfig
	if conf, ok := config[constants.ReportConfigKey]; ok {
		err := mapstructure.Decode(conf, &reportConfig)
		if err != nil {
			cclog.Error(fmt.Sprintf("decode report config failed,%v %v", channelName, conf), err)
		}
	}
	var channelInstance = connectors.ChannelInstance{ConnectorManager: connectorManager, ChannelId: channelId,
		ChannelName: channelName, Scheduler: gocron.NewScheduler(time.Local),
		ChannelCache: connectors.ChannelCache{}, ReportConfig: reportConfig, Connected: false, Enabled: true}
	opc := &OpcUaConnector{ChannelInstance: channelInstance, Config: config}
	opc.Init()
	return opc
}
