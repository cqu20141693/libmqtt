/*
 * Copyright Go-IIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"fmt"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/domain"
	"regexp"
)

func ConnHandler(client mqtt.Client, server string, code byte, err error) {
	if err != nil {
		fmt.Printf("\nconnect to server:%v error:%v\n", server, err)
	} else if code != mqtt.CodeSuccess {
		fmt.Printf("\nconnection rejected by server：%v, code:%v\n", server, code)
	} else {
		fmt.Printf("\nconnected to server:%v\n", server)
	}
	fmt.Printf(domain.LineStart)
}

func PubHandler(client mqtt.Client, topic string, err error) {
	if err != nil {
		fmt.Println("\npub", topic, "failed, error =", err)
	} else {
		fmt.Println("\npub", topic, "success")
	}
	fmt.Printf(domain.LineStart)
}

func SubHandler(client mqtt.Client, topics []*mqtt.Topic, err error) {
	if err != nil {
		fmt.Println("\nsub", topics, "failed, error =", err)
	} else {
		fmt.Println("\nsub", topics, "success")
	}
	fmt.Printf(domain.LineStart)
}

func UnSubHandler(client mqtt.Client, topics []string, err error) {
	if err != nil {
		fmt.Println("\nunsub", topics, "failed, error =", err)
	} else {
		fmt.Println("\nunsub", topics, "success")
	}
	fmt.Printf(domain.LineStart)
}

func NetHandler(client mqtt.Client, server string, err error) {
	fmt.Println("\nconnection to server, error:", err)
	fmt.Printf(domain.LineStart)
	client.Destroy(false)
}
func NetHandlerWithReconnect(client mqtt.Client, server string, err error) {
	fmt.Println("\nconnection to server, error:", err)
	fmt.Printf(domain.LineStart)
	client.Reconnect(server)
}

func TopicHandler(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
	fmt.Println("\n[MSG] topic:", topic, "msg:", string(msg), "qos:", qos)
	fmt.Printf(domain.LineStart)
}

func InvalidQos() {
	fmt.Println("\nqos level should either be 0, 1 or 2")
	fmt.Printf(domain.LineStart)
}

// 配置常量
const (
	WelcomeTopic = "sys/welcome"
	CmdTopic     = "sys/cmd/*"
)

var handlerMap = make(map[*regexp.Regexp]mqtt.TopicHandleFunc, 8)

// AddHandler
// 添加topic 处理器
func AddHandler(topicRegex string, h mqtt.TopicHandleFunc) {
	handlerMap[regexp.MustCompile(topicRegex)] = h
}

func HandleMsg(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
	for r, handleFunc := range handlerMap {
		if r.MatchString(topic) {
			handleFunc(client, topic, qos, msg)
		}
	}
}
func CreatePublishPacket(topic string, qos byte, message string) *mqtt.PublishPacket {
	return &mqtt.PublishPacket{
		TopicName: topic,
		Qos:       qos,
		Payload:   []byte(message),
	}
}
