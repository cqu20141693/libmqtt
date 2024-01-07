package register

import (
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/connectors/mqtt"
	opcua "github.com/goiiot/libmqtt/edge_gateway/connectors/opc"
	"github.com/goiiot/libmqtt/edge_gateway/connectors/simulator"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
)

type InstanceGenerator = func(connManager *connectors.ConnectorManager, channelId string, channelName string, config map[string]interface{}) connectors.Channel

var DefaultConnectors = map[string]InstanceGenerator{
	"Simulator": simulator.NewSimulatorConnector,
	"Mqtt":      mqtt.NewMqttConnector,
	"Opc":       opcua.NewOpcUaConnector,
}

// RegisterConnector 注册协议连接器
//
//	@param protocol 协议
//	@param generator 通道生成器
//	@param override 是否覆盖
func RegisterConnector(protocol string, generator InstanceGenerator, override bool) {
	if _, exist := DefaultConnectors[protocol]; exist {
		if override {
			DefaultConnectors[protocol] = generator
		} else {
			cclog.Warn(fmt.Sprintf("registe connector failed,protocol=%v already exist", protocol))
		}
	} else {
		DefaultConnectors[protocol] = generator
	}

}
