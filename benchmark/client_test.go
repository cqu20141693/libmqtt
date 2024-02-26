/*
 * Copyright Go-IIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package benchmark

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"github.com/google/uuid"
	"math/rand"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	pah "github.com/eclipse/paho.mqtt.golang"

	lib "github.com/goiiot/libmqtt"
)

const (
	testKeepalive         = 3600             // prevent keepalive packet disturb
	testServer            = "localhost:1883" // emqttd broker address
	telemetryMeTopic      = "v1/devices/me/telemetry"
	telemetryTopic        = "v1/gateway/telemetry"
	testBufSize           = 100
	testChild             = true
	testPubCount          = 300
	testConnectCount      = 35
	testConnectChildCount = 50
	timeseriesCount       = 50
	username              = "witeam"
	password              = "witeam@123"
	keepalive             = 60
)

var (
	testTopicMsg = []byte("{\"a\":\"b\",\"c\":\"d\",\"e\":\"f\",\"string\":\"Hello World\"}")
)

// BenchmarkGaPlatformConnect 压测平台连接
//
//	@param b
func BenchmarkGaPlatformConnect(b *testing.B) {
	b.N = testConnectCount
	infos := make([]*domain.MqttClientAddInfo, 0, b.N)
	idPrefix := "benchmark"
	for i := 0; i < b.N; i++ {
		info := domain.NewMqttClientAddInfo(testServer, fmt.Sprintf("%s%d", idPrefix, i), username, password, keepalive)
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
		info := domain.NewMqttClientAddInfo(testServer, gateway, username, password, keepalive)
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
// 压测推送消息
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
		info := domain.NewMqttClientAddInfo(testServer, gateway, username, password, keepalive)
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
					time.Sleep(time.Millisecond * 2000)
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
//	@param productId
//	@param device
//	@return map[string]interface{}
func mockConnectPkt(productId string, device string) map[string]interface{} {
	connPkt := make(map[string]interface{})
	connPkt["productId"] = productId
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

func BenchmarkLibmqttClient(b *testing.B) {
	b.N = testPubCount
	b.ReportAllocs()

	client, err := lib.NewClient(
		lib.WithKeepalive(testKeepalive, 1.2),
		lib.WithCleanSession(true),
	)

	if err != nil {
		b.Error(err)
	}

	_ = client.ConnectServer(testServer,
		lib.WithUnsubHandleFunc(func(client lib.Client, topics []string, err error) {
			if err != nil {
				b.Error(err)
			}
			client.Destroy(true)
		}),
		lib.WithConnHandleFunc(func(client lib.Client, server string, code byte, err error) {
			if err != nil {
				b.Error(err)
			} else if code != lib.CodeSuccess {
				b.Error(code)
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				client.Publish(&lib.PublishPacket{
					TopicName: telemetryMeTopic,
					Payload:   testTopicMsg,
				})
			}
			client.UnSubscribe(telemetryMeTopic)
		}),
	)

	client.Wait()
}

func BenchmarkPahoClient(b *testing.B) {
	// nolint:staticcheck
	b.N = testPubCount
	b.ReportAllocs()

	serverURL, err := url.Parse("tcp://" + testServer)
	if err != nil {
		b.Error(err)
	}

	client := pah.NewClient(&pah.ClientOptions{
		Servers:             []*url.URL{serverURL},
		KeepAlive:           testKeepalive,
		CleanSession:        true,
		ProtocolVersion:     4,
		MessageChannelDepth: testBufSize,
		Store:               pah.NewMemoryStore(),
	})

	t := client.Connect()
	if !t.Wait() {
		b.Fail()
	}

	if err := t.Error(); err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Publish(telemetryMeTopic, 0, false, testTopicMsg)
	}

	t = client.Unsubscribe(telemetryMeTopic)
	if !t.Wait() {
		b.Fail()
	}
	if err := t.Error(); err != nil {
		b.Error(err)
	}

	client.Disconnect(0)
}
