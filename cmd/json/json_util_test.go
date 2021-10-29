package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonStrToMap(t *testing.T) {
	result := &MqttAuthResult{}
	jsonStr := "{\"data\":{\"clientIdentifier\":\"wjfypGlG20211020\",\"username\":\"3TPzNOsG\",\"password\":\"15LymzBG\"},\"code\":\"0000\",\"message\":\"成功\"}"
	bytes := []byte(jsonStr)
	json.Unmarshal(bytes, result)
	fmt.Printf("json 序列化 r=%v \n", result)
}
