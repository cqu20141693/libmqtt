package http

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/goiiot/libmqtt/cmd/json"
	"strings"
)

func GetMirrorMqttInfo(url string) (info *json.MqttInfo, e error) {
	split := strings.Split(url, "?")
	client := resty.New()
	r := client.R()
	if len(split) == 2 {
		r.SetQueryString(split[1]).
			SetHeader("Accept", "application/json")

	}
	response, err := r.Post(split[0])
	if err != nil {
		return nil, err
	}
	jsonStr := fmt.Sprintf(string(response.Body()))
	fmt.Println(jsonStr)
	authResult, err := json.JsonToMqttAuthResult([]byte(jsonStr))
	if err != nil {
		return nil, err
	}
	return &authResult.Data, nil

}
func GetSubMqttInfo(url string, devices []string) (info *json.MqttInfo, e error) {
	// Create a Resty Client
	client := resty.New()
	post, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(devices).
		Post(url)
	if err != nil {
		return
	}
	fmt.Println(string(post.Body()))
	authResult, err := json.JsonToMqttAuthResult(post.Body())
	if err != nil {
		return nil, err
	}
	return &authResult.Data, nil

}
