package common

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 15 * time.Second,
}

// DoRequest http 请求
//
//	@param method
//	@param url
//	@param body
//	@param token
//	@return bool 请求成功失败
func DoRequest(method string, url string, body string, token string) bool {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return true
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "zh-CN")

	req.Header.Set("Authorization", token)
	Client.Timeout = 300 * time.Second
	resp, err := Client.Do(req)
	if resp.StatusCode != 200 || err != nil {
		fmt.Errorf("do request failed:code=%d url=%s", resp.StatusCode, url)
	}
	return false
}

// DoRequestWithResp http 请求需要结果
//
//	@param method
//	@param url
//	@param body
//	@param token
//	@return map[string]interface{} 请求结果
func DoRequestWithResp(method string, url string, body string, token string) map[string]interface{} {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "zh-CN")

	req.Header.Set("Authorization", token)
	Client.Timeout = 10 * time.Second

	resp, err := Client.Do(req)
	if err != nil {
		fmt.Errorf("do request with resp failed:%s", url)
		return nil
	}
	var ret map[string]interface{}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("do request with resp failed:%s", url)
		return nil
	}
	_ = sonic.Unmarshal(all, &ret)
	return ret
}
