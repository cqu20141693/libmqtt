package ga

import (
	"fmt"
	"github.com/goiiot/libmqtt/common"
	"testing"
)

func TestDeviceCategory(t *testing.T) {
	maxNum := 100
	groupName := fmt.Sprintf("直连%d", maxNum)
	AddDeviceCategory(groupName)
	treeMap := FetchGroupTree()
	groupId := treeMap[groupName]
	if groupId != "" {
		BindDeviceGroup(groupId, maxNum-100, maxNum)
	} else {
		_ = fmt.Errorf("group=%s not exist", groupName)
	}
}

func BindDeviceGroup(groupId string, from int, to int) {
	if from >= to {
		_ = fmt.Errorf("from=%d is greater to=%d", from, to)
		return
	}
	var ids string
	for i := from; i < to; i++ {
		ids = ids + fmt.Sprintf("direct%d", i)
		if i < to-1 {
			ids = ids + ","
		}
	}
	bindUrl := address + "/v3/device/group/move"
	config := "{\n  \"ids\": \"" + ids +
		"\",\n  \"targetId\": \"" + groupId +
		"\"\n}"
	// bind
	common.DoRequest("POST", bindUrl, config, token)
}

// AddDeviceCategory
//
//	@param groupId
func AddDeviceCategory(groupName string) {
	deviceGroupUrl := address + "/v3/device/group"

	config := "{\n  \"name\": \"" + groupName +
		"\",\n  \"description\": \"\",\n  \"parentId\": \"\"\n}"
	resp := common.DoRequestWithResp("POST", deviceGroupUrl, config, token)
	fmt.Println(resp)
}

// FetchGroupTree 获取group 关系
//
//	@return map[string]string
func FetchGroupTree() map[string]string {
	treeUrl := address + "/v3/device/group/_query/tree"
	config := "{\"terms\":[" +
		"{\"column\":\"level\",\"value\":1,\"termType\":\"eq\"}," +
		"{\"column\":\"name\",\"value\":\"%直连%\",\"type\":\"and\",\"termType\":\"like\"}" +
		"]}"
	resp := common.DoRequestWithResp("POST", treeUrl, config, token)
	treeMap := make(map[string]string)
	if resp != nil {
		result := resp["result"].([]interface{})

		for _, ret := range result {

			group := ret.(map[string]interface{})
			i := group["name"].(string)
			treeMap[i] = group["id"].(string)
		}
	}
	return treeMap
}
