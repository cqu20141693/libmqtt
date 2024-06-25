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

	TestServer  = "10.168.141.202:30009"
	TestToken   = "Bearer eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxIiwidXNlcl9pZCI6IjEiLCJhenAiOiJlbWJlZC1pYW0iLCJ1bmlxdWVfa2V5IjoiOTMzMTdlZDYtZWVkYi00MGJjLTlmNDItZjM4NTUzYTg2MjQ5IiwiYWNjZXNzX2p0aSI6ImE4M2Q1MDc0LTljZGYtNGJkZC1iMjIyLTZmN2FjYmI4ZTNmZSIsIm5hbWUiOiLotoXnuqfnrqHnkIblkZgiLCJpc3MiOiJodHRwOi8vZ3VjMy1hcGktbW9tLmNhYXMtY2xvdWQtdGVzdC5nZWVnYS5jb20vYXBpL2lhbS8xIiwidHlwIjoiQmVhcmVyIiwicmVhbG0iOiIxIiwibG9naW5fc291cmNlIjoidXNlcm5hbWUtcGFzc3dvcmQiLCJqdGkiOiJhODNkNTA3NC05Y2RmLTRiZGQtYjIyMi02ZjdhY2JiOGUzZmUiLCJpYXQiOjE3MTY3OTUxOTMsImV4cCI6MTcxNjc5Njk5M30.CQyelhbd9iTR9vX6XQ2QINR0soCYC5FCG0FoPCsuI8wnTq7Nm94MGK2CLMgDLCTv-yAyRzhZ04NEneQl0_ouuA"
	TestAddress = "http://iiot-2.caas-cloud-test.gee" +
		"ga.com/api/iot-service"

	EdgeServer  = "172.28.89.214:21883"
	EdgeToken   = "mock guc"
	EdgeAddress = "http://172.28.89.214:18840"
)

// 默认值
var (
	testConnectChildCount = 50
	timeseriesCount       = 50

	//gateway product ID
	gatewayProductId   = "witeamP"
	directDeviceFormat = "direct%d"
	gatewayPrefix      = "benchmark_g"
	gatewayIdFormat    = "benchmark_g%d"

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
func UseCass4Config() {
	fmt.Println("use cass 4 config")
	common.Server = "10.168.141.202:30024"
	address = TestAddress
	token = TestToken
}
