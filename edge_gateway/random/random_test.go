package random

import (
	"encoding/json"
	"fmt"
	"testing"
)

var ProductId = "1739547140905361408"

func TestMockJSON(t *testing.T) {

	devices := 150
	var infos = make([]map[string]interface{}, 0)
	for i := 0; i < devices; i++ {
		//deviceTopic := mockDeviceInfoWithTopic()
		deviceTopic := MockDeviceInfo()

		infos = append(infos, deviceTopic)
	}

	marshal, _ := json.Marshal(infos)
	fmt.Println(string(marshal))
	//cclog.Info(string(marshal))
}

func mockDeviceInfoWithTopic() map[string]interface{} {
	var message = make(map[string]interface{})
	message["topic"] = "v1/gateway/connect"
	message["qos"] = 1

	json := MockDeviceInfo()
	message["message"] = json
	return message
}

func MockDeviceInfo() map[string]interface{} {
	keys := []string{"productName", "deviceId", "deviceName"}
	json := randomJSON(keys)
	json["productId"] = ProductId
	//cclog.Info(json)
	keys = []string{"name", "key"}
	props := 10
	var timeseries = make([]map[string]interface{}, 0)

	for i := 0; i < props; i++ {
		values := randomJSON(keys)
		var types = []string{"Bool", "Integer", "Long", "Double", "String"}
		values["dataType"] = randomByFixName(types)
		timeseries = append(timeseries, values)
	}
	json["timeseries"] = timeseries
	return json
}

func randomByFixName(types []string) string {
	if types == nil || len(types) == 0 {
		return "nil"
	}
	randInt := RandInt(0, len(types))
	return types[randInt]
}

func randomJSON(keys []string) map[string]interface{} {
	var json = make(map[string]interface{})
	for _, key := range keys {
		randString := RandString(6)
		json[key] = randString
	}
	return json
}
