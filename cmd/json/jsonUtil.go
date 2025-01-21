package json

import (
	"encoding/json"
	"github.com/goiiot/libmqtt/common/model"
)

func JsonToMqttAuthResult(jsonBytes []byte) (r *model.MqttAuthResult, e error) {
	result := &model.MqttAuthResult{}
	err := json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
