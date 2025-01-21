package ga

import (
	lib "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/common"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"log"
	"testing"
	"time"
)

func TestEosGatewayConnect(t *testing.T) {

	// 使用第三套环境：配置了mqtt 地址
	UseCass3Config()

	common.MqttConnect(1, func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
		// 连接地址，clientId,username,password
		return common.Server, "1877543064354234368", "witeam", "witeam@123", common.Keepalive, lib.V311, time.Second * 10
	})

	time.Sleep(time.Second * 1000)
}

func TestEosGatewayChildPublishData(t *testing.T) {
	UseCass3Config()
	testConnectChildCount = 20
	timeseriesCount = 5
	pushCount := 10
	duration := time.Millisecond * 1000

	common.MqttPublish(1,
		func(index int) (string, string, string, string, int64, lib.ProtoVersion, time.Duration) {
			return common.Server, "1877543064354234368", directUsername, directPassword, common.Keepalive, lib.V311, time.Second * 10
		},
		func(clientId string) (common.TelemetryFunc, string, int, lib.QosLevel, time.Duration) {

			return mockEosGatewayTelemetryPkt, telemetryTopic, pushCount, lib.QosLevel(1), duration
		})
	log.Println("MqttPublish success")
	time.Sleep(time.Second * 3)
}

func mockEosGatewayTelemetryPkt(gatewayId string) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	children := make([]string, 0, 1)
	children = append(children, "1873999740110884864")
	pkt := mockEosTelemetryPkt(children)
	ret = append(ret, pkt)
	return ret
}

func mockEosTelemetryPkt(children []string) map[string]interface{} {
	telemetryPkt := make(map[string]interface{})
	for _, child := range children {

		device := child
		datas := make([]map[string]interface{}, 0, 1)
		data := make(map[string]interface{})
		data["ts"] = utils.GetTimestamp()
		values := make(map[string]interface{})

		mergeMaps(values, mockFloatValues([]string{"cumulative_run_time1"}))
		mergeMaps(values, mockBoolValues([]string{"status1", "equipment_control_mode1"}))
		data["values"] = values
		datas = append(datas, data)
		telemetryPkt[device] = datas
	}
	return telemetryPkt
}

func mockBoolValues(keys []string) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, key := range keys {
		ret[key] = common.MockBool()
	}
	return ret
}

func mockFloatValues(keys []string) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, key := range keys {
		float := common.MockFloat(10.0, 100000.0)
		ret[key] = float
	}
	return ret
}

func mergeMaps(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}
