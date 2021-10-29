package json

import (
	"encoding/json"
)

func JsonToMqttAuthResult(jsonBytes []byte) (r *MqttAuthResult, e error) {
	result := &MqttAuthResult{}
	err := json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
