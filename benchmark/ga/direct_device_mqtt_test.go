package ga

import (
	"fmt"
	"github.com/goiiot/libmqtt"
	common2 "github.com/goiiot/libmqtt/common"
	"log"
	"testing"
	"time"
)

// BenchmarkGaPlatformGatewayConnect 压测平台连接
//
//	@param b
func BenchmarkDirectPlatformConnect(b *testing.B) {

	b.N = 800
	common2.MqttConnect(b.N, func(index int) (string, string, string, string, int64, libmqtt.ProtoVersion, time.Duration) {
		return common2.Server, fmt.Sprintf(directDeviceFormat, index), directUsername, directPassword, common2.Keepalive, libmqtt.V311, -1
	})
	time.Sleep(time.Second * 5)
}

func BenchmarkDirectPlatformPublish(b *testing.B) {
	//b.N = testDirectConnectCount
	b.N = 800

	common2.MqttPublish(b.N,
		func(index int) (string, string, string, string, int64, libmqtt.ProtoVersion, time.Duration) {
			return common2.Server, fmt.Sprintf(directDeviceFormat, index), directUsername, directPassword, common2.Keepalive, libmqtt.V311, time.Second
		},
		func(clientId string) (common2.TelemetryFunc, string, int, libmqtt.QosLevel, time.Duration) {
			duration := time.Millisecond * 1000
			return mockDirectTelemetryPkt, directTopic, 300, libmqtt.QosLevel(0), duration
		})
}

func mockDirectTelemetryPkt(deviceId string) []map[string]interface{} {
	telemetryPkt := make(map[string]interface{})
	telemetryPkt["deviceId"] = deviceId
	values := make(map[string]interface{})
	mockNum := 100
	for i := 0; i < mockNum; i++ {
		key := fmt.Sprintf("a%d", i)
		values[key] = common2.MockInt(0, 9999)
	}
	telemetryPkt["properties"] = values
	ret := make([]map[string]interface{}, 0, 1)
	ret = append(ret, telemetryPkt)
	return ret
}

func TestMockDirectTelemetryPkt(t *testing.T) {
	pkt := mockDirectTelemetryPkt("test")
	log.Println(pkt)
}
