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
	TestToken   = "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxIiwidXNlcl9pZCI6IjEiLCJhenAiOiJlbWJlZC1pYW0iLCJ1bmlxdWVfa2V5IjoiOWJiZTI5YzQtMmExNy00NjUzLTg4ZjQtNTI5NDQ2NjhkZTdhIiwiYWNjZXNzX2p0aSI6IjEyZDk3YjU2LTU3N2YtNDMyMi04OTNlLTdjZmMxMjhiZWYyNyIsIm5hbWUiOiLotoXnuqfnrqHnkIblkZgiLCJpc3MiOiJodHRwOi8vZ3VjMy1hcGktbW9tLmNhYXMtY2xvdWQtdGVzdC5nZWVnYS5jb20vYXBpL2lhbS8xIiwidHlwIjoiQmVhcmVyIiwicmVhbG0iOiIxIiwibG9naW5fc291cmNlIjoidXNlcm5hbWUtcGFzc3dvcmQiLCJqdGkiOiIxMmQ5N2I1Ni01NzdmLTQzMjItODkzZS03Y2ZjMTI4YmVmMjciLCJpYXQiOjE3MTk0NTM1NDMsImV4cCI6MTcyMDMxNzU0M30.R9ezEEhMiqCOv4oJPP5zclRZft_nbpV9wf8dDCtGJmHOm_RqiTWKbu2OwWWJI42uL1vpsq0F1cKSNNC8Cu_cVA"
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

// UseCass2Config 使用测试环境配置
func UseCass2Config() {
	fmt.Println("use test config")
	common.Server = TestServer
	address = TestAddress
	token = TestToken
}

func UseCass3Config() {
	fmt.Println("use cass 3 config")
	common.Server = "10.168.141.229:33183"
	address = "http://iiot-3.caas-cloud-test.gee" +
		"ga.com/api/iot-service"
	token = TestToken
}

func UseCass4Config() {
	fmt.Println("use cass 4 config")
	common.Server = "10.168.141.202:30024"
	address = "http://iiot-4.caas-cloud-test.gee" +
		"ga.com/api/iot-service"
	token = TestToken
}
