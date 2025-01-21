package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/goiiot/libmqtt/cmd/json"
	"github.com/goiiot/libmqtt/common/model"
)

func GetMirrorMqttInfo(args []string) (info *model.MqttInfo, e error) {
	client := resty.New()
	r := client.R()
	if len(args) == 3 {
		r.SetQueryString("groupKey="+args[0]+"&"+"sn="+args[1]).
			SetHeader("Accept", "application/json")
	} else {
		return
	}
	response, err := r.Post(args[2])
	if err != nil {
		return nil, err
	}
	authResult, err := json.JsonToMqttAuthResult(response.Body())
	if err != nil {
		return nil, err
	}
	return &authResult.Data, nil

}
func GetSubMqttInfo(url string, devices []string) (info *model.MqttInfo, e error) {
	// Create a Resty Client
	client := resty.New()
	post, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(devices).
		Post(url)
	if err != nil {
		return
	}
	authResult, err := json.JsonToMqttAuthResult(post.Body())
	if err != nil {
		return nil, err
	}
	return &authResult.Data, nil

}
