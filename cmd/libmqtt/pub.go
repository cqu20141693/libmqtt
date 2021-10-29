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
	"strconv"
	"strings"

	mqtt "github.com/goiiot/libmqtt"
)

func execPub(client *mqtt.AsyncClient, messages []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	pubs := make([]*mqtt.PublishPacket, 0)
	for _, v := range messages {
		pubStr := strings.Split(v, "#")
		if len(pubStr) != 4 {
			pubUsage()
			return true
		}
		qos, err := strconv.Atoi(pubStr[1])
		if err != nil {
			pubUsage()
			return false
		}
		bytes, err := getBytes(pubStr[2], pubStr[3])
		if err != nil {
			fmt.Printf("topic=%v qos=%v type=%v data=%v encode error \n", pubStr[0], pubStr[1], pubStr[2], pubStr[3])
			continue
		}
		pubs = append(pubs, &mqtt.PublishPacket{
			TopicName: pubStr[0],
			Qos:       mqtt.QosLevel(qos),
			Payload:   bytes,
		})
	}
	client.Publish(pubs...)
	return true
}

func pubUsage() {
	fmt.Println(`p, pub [topic,qos,encodeType,message] [...] clientId - publish topic message(s)`)
}
