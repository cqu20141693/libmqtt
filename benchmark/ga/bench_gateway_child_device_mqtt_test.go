package ga

import (
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	lib "github.com/goiiot/libmqtt"
	mqtt_util "github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/common"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"github.com/google/uuid"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	// 使用测试环境配置
	//UseCass4Config()
	UseCass2Config()
	mqtt_util.AddHandler(mqtt_util.Attributes, attributeHandler)
	common.MqttConnect(1, func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
		return common.Server, fmt.Sprintf("%s%d", "witeamG", index+1), "witeam", "witeam@123", common.Keepalive, lib.V311, time.Second * 10
	})

	time.Sleep(time.Second * 1000)
}

func attributeHandler(client lib.Client, topic string, qos lib.QosLevel, msg []byte) {
	data := make(map[string]interface{})
	reply := make(map[string]interface{})
	_ = json.Unmarshal(msg, &data)
	msgId := data["messageId"]
	reply["messageId"] = msgId
	reply["success"] = true
	reply["timestamp"] = utils.GetTimestamp()
	bytes, _ := json.Marshal(reply)
	// 2.x topic 回复
	client.Publish(&lib.PublishPacket{
		TopicName: "properties/write/reply",
		Qos:       lib.Qos1,
		Payload:   bytes,
	})
}

// BenchmarkGaPlatformGatewayConnect 压测平台连接
//
//	@param b
func BenchmarkGaPlatformGatewayConnect(b *testing.B) {
	// 使用测试环境配置
	UseCass3Config()

	b.N = 10
	common.MqttConnect(b.N, func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
		return common.Server, fmt.Sprintf("%s%d", gatewayPrefix, index), common.Username, common.Password, common.Keepalive, lib.V311, time.Second * 10
	})

	time.Sleep(time.Second * 10)
}

// BenchmarkGaPlatformConnectChild
// 通过connect 报文创建子设备和物模型
//
//	@param b
func BenchmarkGaPlatformConnectChild(b *testing.B) {
	UseCass2Config()
	testConnectChildCount = 20
	timeseriesCount = 50
	b.N = 10
	common.MqttPublish(b.N,
		func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
			return common.Server, fmt.Sprintf("%s%d", gatewayPrefix, index), directUsername, directPassword, common.Keepalive, lib.V311, time.Second * 10
		},
		func(clientId string) (common.TelemetryFunc, string, int, lib.QosLevel, time.Duration) {
			duration := time.Millisecond * 3000
			return mockGatewayDeviceTSLPkt, // 报文生成器
				"v1/gateway/connect", //topic
				1, // 发送次数
				lib.QosLevel(1), //qos
				duration // 发送延迟
		})
	log.Println("MqttPublish success")
	time.Sleep(time.Second * 30)

}

// BenchmarkEdgeRuleService
// 压测edge rule 服务 2000 点位
//
//	@param b
func BenchmarkEdgeRuleService(b *testing.B) {
	common.Server = EdgeServer
	address = EdgeAddress
	token = EdgeToken
	testConnectChildCount = 20
	timeseriesCount = 50
	b.N = 1
	pushCount := 18000
	// 每秒10条
	duration := time.Millisecond * 100

	common.MqttPublish(b.N,
		func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
			return common.Server, fmt.Sprintf("%s%d", gatewayPrefix, index), common.Username, common.Password,
				common.Keepalive, lib.V311, time.Second * 3
		},
		func(clientId string) (common.TelemetryFunc, string, int, lib.QosLevel, time.Duration) {

			return mockGatewayTelemetryPkt, telemetryTopic, pushCount, lib.QosLevel(1), duration
		})
	log.Println("MqttPublish success")
	time.Sleep(time.Second * 3)
}

// BenchmarkGaPlatformPublish
// 压测子设备推送消息
//
//	@param b
func BenchmarkGaPlatformPublish(b *testing.B) {
	UseCass2Config()
	testConnectChildCount = 20
	timeseriesCount = 5
	b.N = 10
	pushCount := 1000
	duration := time.Millisecond * 1000

	common.MqttPublish(b.N,
		func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
			return common.Server, fmt.Sprintf("%s%d", gatewayPrefix, index), directUsername, directPassword, common.Keepalive, lib.V311, time.Second * 10
		},
		func(clientId string) (common.TelemetryFunc, string, int, lib.QosLevel, time.Duration) {

			return mockGatewayTelemetryPkt, telemetryTopic, pushCount, lib.QosLevel(1), duration
		})
	log.Println("MqttPublish success")
	time.Sleep(time.Second * 3)

}

func BenchmarkGaPlatformGatewayPublish(b *testing.B) {
	b.N = 12
	common.MqttPublish(b.N,
		func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
			return common.Server, fmt.Sprintf("%s%d", gatewayPrefix, index), directUsername, directPassword, common.Keepalive, lib.V311, time.Second * 10
		},
		func(clientId string) (common.TelemetryFunc, string, int, lib.QosLevel, time.Duration) {
			duration := time.Millisecond * 1000
			return mockGatewayMeTelemetryPkt, telemetryMeTopic, 1, lib.QosLevel(1), duration
		})
	log.Println("MqttPublish success")
	time.Sleep(time.Second * 3)

}

// mockConnectPkt
// 模拟连接topic报文
//
//	@param productId
//	@param device
//	@return map[string]interface{}
func mockConnectPkt(productId string, device string, propsCount int) map[string]interface{} {
	connPkt := make(map[string]interface{})
	connPkt["productId"] = productId
	connPkt["deviceId"] = device
	connPkt["channelId"] = uuid.New().String()
	connPkt["timeseries"] = mockTimeseries(propsCount)
	return connPkt
}

// mockTimeseries
// 模拟属性模型
//
//	@return []map[string]interface{}
func mockTimeseries(propsCount int) []map[string]interface{} {
	timeseries := make([]map[string]interface{}, 0, propsCount)
	for i := 0; i < propsCount; i++ {
		property := make(map[string]interface{})
		property["name"] = "a" + strconv.Itoa(i)
		property["key"] = "a" + strconv.Itoa(i)
		property["dataType"] = randomDataType()
		timeseries = append(timeseries, property)
	}
	return timeseries
}

func randomDataType() string {

	switch common.MockInt(0, 10) {
	case 0:
		return "Integer"
	case 1:
		return "Long"
	case 3:
		return "Double"
	case 4:
		return "Float"
	case 5:
		return "String"
	case 6:
		return "Bool"
	default:
		return "String"
	}

}

// mockTelemetryPkt
// 模拟遥测报文
//
//	@param childPrefix
//	@return map[string]interface{}
func mockTelemetryPkt(childPrefix string) map[string]interface{} {

	telemetryPkt := make(map[string]interface{})
	mockNum := testConnectChildCount
	for i := 0; i < mockNum; i++ {
		device := childPrefix + strconv.Itoa(i)
		datas := make([]map[string]interface{}, 0, 1)
		data := make(map[string]interface{})
		data["ts"] = utils.GetTimestamp()
		data["values"] = mockValues()
		datas = append(datas, data)
		telemetryPkt[device] = datas
	}
	return telemetryPkt
}

// mockValues
// 模拟属性数据
//
//	@return map[string]interface{}
func mockValues() map[string]interface{} {
	values := make(map[string]interface{})
	for i := 0; i < timeseriesCount; i++ {
		key := "a" + strconv.Itoa(i)
		//values[key] = random.RandInt(10, 10000)
		values[key] = common.MockInt(10, 10000)
	}
	return values
}

// mockGatewayTelemetryPkt 模拟网关子设备遥测数据
//
//	@param gateWayId
//	@return []map[string]interface{}
func mockGatewayTelemetryPkt(gatewayId string) []map[string]interface{} {
	index := gatewayId[len(gatewayPrefix):]
	childPrefix := "benchmarkCC_" + index + "_"
	ret := make([]map[string]interface{}, 0, 1)
	pkt := mockTelemetryPkt(childPrefix)
	ret = append(ret, pkt)
	return ret
}

func TestMockGatewayTelemetryPkt(t *testing.T) {
	pkt := mockGatewayTelemetryPkt("benchmark6")
	marshal, _ := sonic.Marshal(pkt)
	log.Println(string(marshal))
}

// mockGatewayDeviceTSLPkt
//
//	@param deviceId 网关id
//	@return map[string]interface{}
func mockGatewayDeviceTSLPkt(deviceId string) []map[string]interface{} {
	index := deviceId[len(gatewayPrefix):]
	productId := productIdPrefix + "0"
	childPrefix := "benchmarkCC_" + index + "_"

	//timeseriesCount = 10

	ret := make([]map[string]interface{}, 0, testConnectChildCount)
	for i := 0; i < testConnectChildCount; i++ {
		child := childPrefix + strconv.Itoa(i)
		connPkt := mockConnectPkt(productId, child, timeseriesCount)
		marshal, _ := sonic.Marshal(connPkt)
		log.Println(string(marshal))
		ret = append(ret, connPkt)
	}

	return ret
}

func TestMockGatewayDeviceTSLPkt(t *testing.T) {
	pkt := mockGatewayDeviceTSLPkt(gatewayPrefix + "0")
	marshal, _ := sonic.Marshal(pkt)

	log.Println(string(marshal))
}

func mockGatewayMeTelemetryPkt(gatewayId string) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	values := make(map[string]interface{})
	mockNum := 100
	for i := 0; i < mockNum; i++ {
		key := fmt.Sprintf("a%d", i)
		values[key] = common.MockInt(0, 9999)
	}
	ret = append(ret, values)
	return ret
}
