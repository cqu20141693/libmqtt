package modbus

import (
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/simonvetter/modbus"
)

type TcpConnector struct {
	connectors.ChannelInstance
	Config map[string]interface{} // 整体配置
	// 连接配置，根据协议不同，连接配置不同
	// 每一种配置可以定义自己的struct
	ConnectConfig map[string]interface{}
	Client        *modbus.ModbusClient
}

func (t TcpConnector) Init() {

	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) Open() {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) Close(remote bool) {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) IsConnected() bool {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) GetId() string {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) ServerSideRpcHandler(eventData connectors.RpcEvent) {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) ServerSideRpcHandlerSync(eventData connectors.RpcEvent) connectors.RpcResult {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) AttributesUpdate(eventData connectors.UpdateTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) AttributesUpdateSync(eventData connectors.UpdateTagEvent) connectors.UpdateTagResult {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) AttributesRead(eventData connectors.ReadTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (t TcpConnector) AttributesReadSync(eventData connectors.ReadTagEvent) connectors.ReadTagResult {
	//TODO implement me
	panic("implement me")
}
