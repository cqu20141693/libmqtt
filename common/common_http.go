package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 15 * time.Second,
}

func DoRequestWithTimeout(method string, url string, body string, token string, ctx context.Context) bool {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "zh-CN")

	req.Header.Set("Authorization", token)
	Client.Timeout = 300 * time.Second
	resp, err := Client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println(fmt.Sprintf("do request failed:err=%v", err))
		return false
	} else if resp.StatusCode != 200 {
		all, _ := io.ReadAll(resp.Body)
		fmt.Println(fmt.Sprintf("do request failed:code=%d reson=%s url=%s", resp.StatusCode, string(all), url))
		return false
	}
	return true
}

// DoRequest api 请求
//
//	@param method
//	@param url
//	@param body
//	@param token
//	@return bool false请求成功失败,ture 失败
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
		fmt.Println(fmt.Sprintf("do request failed:code=%d url=%s", resp.StatusCode, url))
	}
	return false
}

// DoRequestWithResp api 请求需要结果
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
		_ = fmt.Errorf("do request with resp failed:%s", url)
		return nil
	}
	var ret map[string]interface{}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = fmt.Errorf("do request with resp failed:%s", url)
		return nil
	}
	_ = sonic.Unmarshal(all, &ret)
	return ret
}
