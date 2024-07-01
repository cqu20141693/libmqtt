package common

import (
	"fmt"
	lib "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/common"
	"sync/atomic"
	"testing"
	"time"
)

func TestReconnect(t *testing.T) {
	common.Server = "localhost:2883"
	connInfoFunc := func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
		return common.Server, fmt.Sprintf("%s%d", "witeamG", index+1), "witeam", "witeam@123", common.Keepalive, lib.V311, time.Second * 10
	}
	common.MqttConnectWithLog(1, connInfoFunc, lib.Info)
	time.Sleep(time.Second * 100)
}

func TestRecoverData(t *testing.T) {

	common.Server = "localhost:2883"
	connInfoFunc := func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
		return common.Server, fmt.Sprintf("%s%d", "demo-", index+1), "witeam", "witeam@123", common.Keepalive, lib.V311, time.Second * 3
	}
	pubInfoFunc := func(clientId string) (common.TelemetryFunc, string, int, lib.QosLevel, time.Duration) {
		duration := time.Millisecond * 1000
		return mockDirectTelemetryPkt, "/report-property", 60, lib.QosLevel(0), duration
	}
	common.MqttPublish(1,
		connInfoFunc,
		pubInfoFunc)
}

var nextID uint32

func mockDirectTelemetryPkt(deviceId string) []map[string]interface{} {
	telemetryPkt := make(map[string]interface{})
	telemetryPkt["deviceId"] = deviceId
	values := make(map[string]interface{})

	values["a"] = atomic.AddUint32(&nextID, 1)

	telemetryPkt["properties"] = values
	ret := make([]map[string]interface{}, 0, 1)
	ret = append(ret, telemetryPkt)
	return ret
}
