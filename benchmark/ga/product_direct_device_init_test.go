package ga

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/goiiot/libmqtt/common"
	"strconv"
	"testing"
)

func initDirectConfig() {
	productIdPrefix = "directP"
	protocol = "demo-protocol"
	deviceType = "device"
	protocolName = "demo-protocol"
}

func TestMockDirectProduct(t *testing.T) {
	initDirectConfig()
	MockDirectProduct(1, address, token)

}

func TestMockDirectDevice(t *testing.T) {
	initDirectConfig()

	MockDevice(productIdPrefix+"0", 500, 800, token, address)
}

// MockDirectProduct 创建直连产品
//
//	@param count
//	@param address
//	@param token
func MockDirectProduct(count int, address string, token string) {

	deployProductFormat := address + "/device/product/%s/%s/deploy"

	createProduct := address + "/device/product" //PATCH

	createPBodyFormat := "{\"id\":\"%s\"," +
		"\"name\":\"%s\"," +
		"\"version\":\"1.0.0\"," +
		"\"messageProtocol\":\"" + protocol +
		"\",\"transportProtocol\":\"MQTT\"," +
		"\"deviceType\":\"" + deviceType +
		"\",\"describe\":\"\",\"extendConfig\":{}," +
		"\"classifiedName\":\"\",\"protocolName\":\"" + protocolName +
		"\",\"metadata\":\"{}\"}"
	for i := 0; i < count; i++ {
		id := productIdPrefix + strconv.Itoa(i)
		body := fmt.Sprintf(createPBodyFormat, id, id)

		// create product
		common.DoRequest("PATCH", createProduct, body, token)

		// query  versions
		productVersionFormat := address + "/product/version/%s/versions"
		url := fmt.Sprintf(productVersionFormat, id)
		resp := common.DoRequestWithResp("GET", url, "", token)
		if resp != nil {
			result := resp["result"].([]interface{})
			var versionId string
			for _, ret := range result {
				version := ret.(map[string]interface{})
				versionId = version["id"].(string)
			}
			if versionId != "" {
				// auth config
				authConfigFormat := address + "/product/version/%s"
				config := "{\"configuration\": {\n    \"Username\": \"witeam\",\n    \"Password\": \"witeam@123\"\n  }\n}"
				common.DoRequest("PUT", fmt.Sprintf(authConfigFormat, versionId), config, token)
				// deploy p
				common.DoRequest("POST", fmt.Sprintf(deployProductFormat, id, versionId), "{}", token)
			}
		}
	}

}

func MockProductTSL(productId string) {

	versionId := GetProductVersionId(productId)
	if versionId == "" {
		_ = fmt.Errorf("GetProductVersionId failed:%s", productId)
		return
	}
	tslFormat := "{\"metadata\": \"{\\\"tags\\\":[{\\\"id\\\":\\\"Location\\\",\\\"name\\\":\\\"位置信息\\\",\\\"valueType\\\":{\\\"id\\\":\\\"geoPoint\\\",\\\"latProperty\\\":\\\"lat\\\",\\\"lonProperty\\\":\\\"lon\\\",\\\"name\\\":\\\"地理位置\\\",\\\"type\\\":\\\"geoPoint\\\"}}]," +
		"\\\"properties\\\":[{\\\"id\\\":\\\"a0\\\",\\\"name\\\":\\\"a0\\\",\\\"description\\\":\\\"\\\",\\\"valueType\\\":{\\\"type\\\":\\\"int\\\"},\\\"expands\\\":{\\\"readOnly\\\":[\\\"report\\\"],\\\"propertyType\\\":\\\"timeseries\\\",\\\"source\\\":\\\"device\\\",\\\"storageType\\\":\\\"direct\\\"}}]}\"}"
	tslUrlFormat := address + "/product/version/%s"
	fmt.Println(tslFormat, tslUrlFormat)
}

func TestTslMock(t *testing.T) {
	tslFormat := `{
  "metadata": "{\"tags\":[{\"id\":\"Location\",\"name\":\"位置信息\",\"valueType\":{\"id\":\"geoPoint\",\"latProperty\":\"lat\",\"lonProperty\":\"lon\",\"name\":\"地理位置\",\"type\":\"geoPoint\"}}],\"properties\":[{\"id\":\"a0\",\"name\":\"a0\",\"description\":\"\",\"valueType\":{\"type\":\"int\"},\"expands\":{\"readOnly\":[\"report\"],\"propertyType\":\"timeseries\",\"source\":\"device\",\"storageType\":\"direct\"}}]}"
}`
	template := make(map[string]interface{})
	err := sonic.UnmarshalString(tslFormat, template)
	if err != nil {
		return
	}
	metadata := template["metadata"].(string)
	tsl := make(map[string]interface{})
	err = sonic.UnmarshalString(metadata, tsl)
	if err != nil {
		return
	}
	fmt.Println(tsl["tags"])

}

// MockDevice 创建设备(可以支持直连，网关)
//
//	@param productId
//	@param from 开始id
//	@param to 结束id
//	@param token
//	@param address
func MockDevice(productId string, from int, to int, token string, address string) {
	if from >= to {
		_ = fmt.Errorf("MockDevice from=%d is greater than to=%d", from, to)
		return
	}
	for i := from; i < to; i++ {

		deviceId := fmt.Sprintf(directDeviceFormat, i)
		createUrl := address + "/device/instance"
		deployUrl := address + fmt.Sprintf("/device/instance/%s/deploy", deviceId)

		versionId := GetProductVersionId(productId)
		if versionId == "" {
			_ = fmt.Errorf("GetProductVersionId failed:%s", productId)
			return
		}

		temp := fmt.Sprintf("{\n  \"id\": \"%s\","+
			"\n  \"name\": \"%s\",\n"+
			"  \"productId\": \""+productId+
			"\",\n  \"version\": \""+versionId+
			"\",\n  \"describe\": \"\",\n"+
			"  \"productName\": \"压测\",\n"+
			"  \"description\": \"\"\n}", deviceId, deviceId)

		method := "POST"

		done := common.DoRequest(method, createUrl, temp, token)
		if done {
			return
		}
		done = common.DoRequest(method, deployUrl, "{}", token)
		if done {
			return
		}

	}

}
