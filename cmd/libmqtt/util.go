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

package main

import (
	"fmt"
	mqtt "github.com/goiiot/libmqtt"
)

func connHandler(client mqtt.Client, server string, code byte, err error) {
	if err != nil {
		fmt.Println("\nconnect to server error:", err)
	} else if code != mqtt.CodeSuccess {
		fmt.Println("\nconnection rejected by server, code:", code)
	} else {
		fmt.Println("\nconnected to server")
	}
	fmt.Printf(lineStart)
}

func pubHandler(client mqtt.Client, topic string, err error) {
	if err != nil {
		fmt.Println("\npub", topic, "failed, error =", err)
	} else {
		fmt.Println("\npub", topic, "success")
	}
	fmt.Printf(lineStart)
}

func subHandler(client mqtt.Client, topics []*mqtt.Topic, err error) {
	if err != nil {
		fmt.Println("\nsub", topics, "failed, error =", err)
	} else {
		fmt.Println("\nsub", topics, "success")
	}
	fmt.Printf(lineStart)
}

func unSubHandler(client mqtt.Client, topics []string, err error) {
	if err != nil {
		fmt.Println("\nunsub", topics, "failed, error =", err)
	} else {
		fmt.Println("\nunsub", topics, "success")
	}
	fmt.Printf(lineStart)
}

func netHandler(client mqtt.Client, server string, err error) {
	fmt.Println("\nconnection to server, error:", err)
	fmt.Printf(lineStart)
}

func topicHandler(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
	fmt.Println("\n[MSG] topic:", topic, "msg:", string(msg), "qos:", qos)
	fmt.Printf(lineStart)
}

func invalidQos() {
	fmt.Println("\nqos level should either be 0, 1 or 2")
	fmt.Printf(lineStart)
}
