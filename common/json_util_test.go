package common

import (
	"encoding/json"
	"fmt"
	model "github.com/goiiot/libmqtt/common/model"
	"testing"
)

func TestJsonStrToMap(t *testing.T) {
	result := &model.MqttAuthResult{}
	jsonStr := "{\"data\":{\"clientIdentifier\":\"wjfypGlG20211020\",\"username\":\"3TPzNOsG\",\"password\":\"15LymzBG\"},\"code\":\"0000\",\"message\":\"成功\"}"
	bytes := []byte(jsonStr)
	err := json.Unmarshal(bytes, result)
	if err != nil {
		return
	}
	fmt.Printf("json 序列化 r=%v \n", result)
}
