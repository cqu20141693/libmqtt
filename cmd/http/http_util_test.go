package http

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

func TestGetMqttInfo(t *testing.T) {
	info, err := GetMirrorMqttInfo([]string{"tPH6EZy6UbIxHxkg", "dVwEOrHihSRHtZTm", "http://172.30.203.22:11093/api/device/authenticate/config/getMirrorAuth"})
	if err != nil {
		fmt.Println("get mqtt info error")
	}
	fmt.Printf("mqtt info =%v", info)
}

func TestGetMqttInfo2(t *testing.T) {
	devices := []string{"dVwEOrHihSRHtZTm", "MQZddYPwBBepuXBx", "WVZhKiRnDDzTyljo"}
	url := "http://172.30.203.22:11093/api/device/authenticate/config/getSubLoginAuthByGroupKey?groupKey=tPH6EZy6UbIxHxkg&appKey=0sxsPsb5lBCEwXa0"
	info, err := GetSubMqttInfo(url, devices)
	if err != nil {
		fmt.Println("get mqtt info error")
	}
	fmt.Printf("mqtt info =%v", info)
}
func TestGet(t *testing.T) {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get("http://172.30.203.22:11093/api/device/authenticate/config/getMirrorAuth?groupKey=tPH6EZy6UbIxHxkg&sn=dVwEOrHihSRHtZTm")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", string(resp.Body()))
	fmt.Println("  Resp       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

func TestPost(t *testing.T) {
	// Create a Resty Client
	client := resty.New()

	post, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]string{"dVwEOrHihSRHtZTm", "MQZddYPwBBepuXBx", "WVZhKiRnDDzTyljo"}).
		Post("http://172.30.203.22:11093/api/device/authenticate/config/getSubLoginAuthByGroupKey?groupKey=tPH6EZy6UbIxHxkg&appKey=0sxsPsb5lBCEwXa0")
	if err != nil {
		return
	}

	fmt.Println(string(post.Body()))
}
