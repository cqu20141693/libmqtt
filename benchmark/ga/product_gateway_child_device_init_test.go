package ga

import (
	"fmt"
	"github.com/goiiot/libmqtt/common"
	"strconv"
	"testing"
)

// TestMockGateway
// 创建网关设备
//
//	@param t

func TestMockGateway(t *testing.T) {
	//gatewayProductId = "v4"
	UseTestConfig()
	// 网关产品可以直接创建，还需要区分数采能力
	MockGatewayDevice(10, token, address)
}

// MockGatewayDevice
// 模拟网关设备
//
//	@param count
//	@param token
//	@param address
func MockGatewayDevice(count int, token string, address string) {

	for i := 0; i < count; i++ {

		deviceId := fmt.Sprintf(gatewayIdFormat, i)
		createUrl := address + "/device/instance"
		deployUrl := address + fmt.Sprintf("/device/instance/%s/deploy", deviceId)

		versionId := GetProductVersionId(gatewayProductId)
		if versionId == "" {
			fmt.Println(fmt.Sprintf("GetProductVersionId failed:%s", gatewayProductId))
		}

		temp := fmt.Sprintf("{\n  \"id\": \"%s\","+
			"\n  \"name\": \"%s\",\n"+
			"  \"productId\": \""+gatewayProductId+
			"\",\n  \"version\": \""+versionId+
			"\",\n  \"describe\": \"\",\n"+
			"  \"productName\": \"压测\",\n"+
			"  \"description\": \"\"\n}", deviceId, deviceId)

		method := "POST"

		failed := common.DoRequest(method, createUrl, temp, token)
		if failed {
			return
		}
		failed = common.DoRequest(method, deployUrl, "{}", token)
		if failed {
			return
		}

	}

}

// TestMockChildDevice
// 创建子设备
//
//	@param t
func TestMockChildDevice(t *testing.T) {

	UseTestConfig()
	// mock child product
	MockProduct(1, address, token)

	// mock device,
	// 子设备通过connect报文自动创建设备和物模型

}

func MockProduct(count int, address string, token string) {

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
		// create p

		common.DoRequest("PATCH", createProduct, body, token)

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
			//deploy p
			common.DoRequest("POST", fmt.Sprintf(deployProductFormat, id, versionId), "{}", token)
		}
	}

}

// GetProductVersionId
// 根据产品id 获取versionId
//
//	@param id
//	@return string
func GetProductVersionId(id string) string {

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
		return versionId
	}
	return ""
}

// TestBatchDelProduct
// 批量删除全量产品
//
//	@param t
func TestBatchDelProduct(t *testing.T) {

	pageQuery := address + "/v3/device/product/queryPager"
	body := "{\"pageSize\":500,\"currentPage\":0,\"sorts\":[{\"name\":\"createTime\",\"order\":\"desc\"}],\"terms\":[{\"column\":\"gatewayProductId\",\"value\":\"0\",\"termType\":\"dev-product-type\"}]}"

	method := "POST"

	ret := common.DoRequestWithResp(method, pageQuery, body, token)
	if ret != nil {
		result := ret["result"].(map[string]interface{})
		records := result["records"].([]interface{})
		for _, record := range records {
			m := record.(map[string]interface{})

			if m["deviceCount"].(float64) == 0 {
				fmt.Println(fmt.Sprintf("%s deviceCount=0", m["name"]))
				delUrl := address + fmt.Sprintf("/product/version/%s", m["currVersion"])
				if m["state"].(float64) == 0 {

					common.DoRequest("DELETE", delUrl, "", token)
				} else {
					unDeployUrl := address + fmt.Sprintf("/device/product/%s/%s/undeploy", m["id"], m["currVersion"])

					common.DoRequest("POST", unDeployUrl, "{}", token)
					common.DoRequest("DELETE", delUrl, "", token)
				}
			}
		}
		fmt.Println(len(ret))
	}
}
