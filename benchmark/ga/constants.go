package ga

import (
	"fmt"
	"github.com/goiiot/libmqtt/common"
)

const (
	telemetryMeTopic = "v1/devices/me/telemetry"
	telemetryTopic   = "v1/gateway/telemetry"
	directTopic      = "/report-property"
	LocalToken       = "tt"
	LocalAddress     = "http://localhost:8840"

	TestServer  = "localhost:1883"
	TestToken   = "Bearer eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIyMjMzMjA2MTY5NTgwMTM0NDAiLCJ1c2VyX2lkIjoiMjIzMzIwNjE2OTU4MDEzNDQwIiwiYXpwIjoiZW1iZWQtaWFtIiwidW5pcXVlX2tleSI6IjVmYzlkNzQxLWQzNTgtNDY0Yy04NzEwLWMwNTMwNWQ1NDA5OSIsImFjY2Vzc19qdGkiOiI0NjdkNjUyNi0wZjE4LTRkNGItOWJmZC0yYjIyMzE0YThiMTUiLCJuYW1lIjoiSUlPVOa1i-ivlS1ORVciLCJpc3MiOiJodHRwOi8vZ3VjMy1hcGktdGVzdC5nZWVnYS5jb20vYXBpL2lhbS8xIiwidHlwIjoiQmVhcmVyIiwicmVhbG0iOiIxIiwibG9naW5fc291cmNlIjoibW9iaWxlLXBhc3N3b3JkIiwianRpIjoiNDY3ZDY1MjYtMGYxOC00ZDRiLTliZmQtMmIyMjMxNGE4YjE1IiwiaWF0IjoxNzA5NjIzMTY2LCJleHAiOjE3MDk4ODIzNjZ9.ao9Ntq0wM-MKARp39MxRfLn3qhm2OaNNFyuaizfabcgU9bxhsj_IUus7q_5FrvETmMIrd-HgweH-dI9w1FZ4vg"
	TestAddress = "https://iiot-3-test.ge" +
		"ega.com/api/iot-service"
)

// 默认值
var (
	testConnectChildCount = 50
	timeseriesCount       = 50

	//gateway product ID
	gatewayProductId   = "witeamP"
	directDeviceFormat = "direct%d"
	gatewayPrefix      = "benchmark"

	address         = LocalAddress
	token           = LocalToken
	protocol        = "g" + "aea-protocol-v1"
	productIdPrefix = "benchmarkCP"
	deviceType      = "childrenDevice"
	protocolName    = "广域工业数采网关"

	directUsername = common.Username
	directPassword = common.Password
)

// UseTestConfig 使用测试环境配置
func UseTestConfig() {
	fmt.Println("use test config")
	common.Server = TestServer
	address = TestAddress
	token = TestToken
}
