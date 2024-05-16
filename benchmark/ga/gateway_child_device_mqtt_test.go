package ga

import (
	"fmt"
	"github.com/bytedance/sonic"
	lib "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/common"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	testTopicMsg = []byte("{\"a\":\"b\",\"c\":\"d\",\"e\":\"f\",\"string\":\"Hello World\"}")
)

// BenchmarkGaPlatformConnect 压测平台连接
//
//	@param b
func BenchmarkGaPlatformConnect(b *testing.B) {
	// 使用测试环境配置
	UseTestConfig()

	b.N = testConnectCount

	infos := make([]*domain.MqttClientAddInfo, 0, b.N)
	idPrefix := "benchmark"
	for i := 0; i < b.N; i++ {
		info := domain.NewMqttClientAddInfo(common.Server, fmt.Sprintf("%s%d", idPrefix, i), common.Username, common.Password, common.Keepalive)
		infos = append(infos, info)

		go func(addInfo *domain.MqttClientAddInfo) {
			client, err := mqtt.CreatClient(addInfo)
			if err != nil {
				_ = fmt.Errorf("creat g platform mqtt client failed, %v", info)
			}
			client.Wait()
		}(info)

	}

	time.Sleep(time.Second * 10)
}

// BenchmarkGaPlatformConnectChild
// 通过connect 报文创建子设备和物模型
//
//	@param b
func BenchmarkGaPlatformConnectChild(b *testing.B) {
	b.N = testConnectCount
	infos := make([]*domain.MqttClientAddInfo, 0, b.N)
	idPrefix := "benchmark"
	group := sync.WaitGroup{}
	group.Add(b.N)

	for i := 0; i < b.N; i++ {
		productId := "benchmarkP" + strconv.Itoa(i)
		childPrefix := "benchmarkC" + strconv.Itoa(i) + "_"
		gateway := fmt.Sprintf("%s%d", idPrefix, i)
		info := domain.NewMqttClientAddInfo(common.Server, gateway, common.Username, common.Password, common.Keepalive)
		infos = append(infos, info)

		go func(addInfo *domain.MqttClientAddInfo) {
			client, err := mqtt.CreatClient(addInfo)
			if err != nil {
				_ = fmt.Errorf("creat g platform mqtt client failed, %v", info)
			}
			time.Sleep(time.Second * 3)
			for i := 0; i < testConnectChildCount; i++ {
				child := childPrefix + strconv.Itoa(i)
				connPkt := mockConnectPkt(productId, child)
				connPktBytes, err := sonic.Marshal(connPkt)
				if err != nil {
					_ = fmt.Errorf("sonic marchal faild: %v", err)
					continue
				}
				// publish connect
				client.Publish(&lib.PublishPacket{
					TopicName: "v1/gateway/connect",
					Payload:   connPktBytes,
				})
			}

			group.Done()
		}(info)

	}
	group.Wait()
}

// BenchmarkGaPlatformPublish
// 压测子设备推送消息
//
//	@param b
func BenchmarkGaPlatformPublish(b *testing.B) {
	b.N = testConnectCount
	infos := make([]*domain.MqttClientAddInfo, 0, b.N)
	idPrefix := "benchmark"
	group := sync.WaitGroup{}
	group.Add(b.N)
	for i := 0; i < b.N; i++ {

		childPrefix := "benchmarkC" + strconv.Itoa(i) + "_"
		gateway := fmt.Sprintf("%s%d", idPrefix, i)
		info := domain.NewMqttClientAddInfo(common.Server, gateway, common.Username, common.Password, common.Keepalive)
		infos = append(infos, info)

		go func(addInfo *domain.MqttClientAddInfo) {
			client, err := mqtt.CreatClient(addInfo)
			if err != nil {
				_ = fmt.Errorf("creat g platform mqtt client failed, %v", info)
			}
			time.Sleep(time.Second * 3)
			b.N = testPubCount
			for i := 0; i < b.N; i++ {
				if testChild {
					time.Sleep(time.Millisecond * 1000)
					telemetryPkt := mockTelemetryPkt(childPrefix)
					telemetryPktBytes, err := sonic.Marshal(telemetryPkt)
					if err != nil {
						_ = fmt.Errorf("sonic marchal faild: %v", err)
						continue
					}
					client.Publish(&lib.PublishPacket{
						TopicName: telemetryTopic,
						Payload:   telemetryPktBytes,
					})
				} else {
					// 网关自上数
					time.Sleep(time.Millisecond * 10)
					client.Publish(&lib.PublishPacket{
						TopicName: telemetryMeTopic,
						Payload:   testTopicMsg,
					})
				}

			}
			group.Done()
		}(info)

	}
	group.Wait()
	time.Sleep(time.Second * testPubCount)
}

// mockConnectPkt
// 模拟连接topic报文
//
//	@param gatewayProductId
//	@param device
//	@return map[string]interface{}
func mockConnectPkt(productId string, device string) map[string]interface{} {
	connPkt := make(map[string]interface{})
	connPkt["gatewayProductId"] = productId
	connPkt["deviceId"] = device
	connPkt["channelId"] = uuid.New().String()
	connPkt["timeseries"] = mockTimeseries()
	return connPkt
}

// mockTimeseries
// 模拟属性模型
//
//	@return []map[string]interface{}
func mockTimeseries() []map[string]interface{} {
	timeseries := make([]map[string]interface{}, 0, timeseriesCount)
	for i := 0; i < timeseriesCount; i++ {
		property := make(map[string]interface{})
		property["name"] = "a" + strconv.Itoa(i)
		property["key"] = "a" + strconv.Itoa(i)
		property["dataType"] = "String"
		timeseries = append(timeseries, property)
	}
	return timeseries
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
		values[key] = mockInt(10, 10000)
	}
	return values
}

// mockInt
// 模拟整数
//
//	@param min
//	@param max
//	@return int
func mockInt(min int, max int) int {
	return min + (rand.Int() % (max - min))
}
