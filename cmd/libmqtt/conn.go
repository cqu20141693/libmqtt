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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/goiiot/libmqtt/cmd/crypto"
	"github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/domain"
	"os"
	"strconv"
	"strings"

	mqtt "github.com/goiiot/libmqtt"
)

func execConn(args []string, version mqtt.ProtoVersion, info *domain.ClientInfo) (client mqtt.Client, err error) {
	if len(args) < 5 {
		return nil, err
	}

	if len(strings.Split(args[0], ":")) != 2 {
		return nil, err
	}

	options := make([]mqtt.Option, 0)
	if len(args) == 6 {
		opts := strings.Split(args[1], ",")
		var willQos mqtt.QosLevel
		var sslSkipVerify, ssl, will, willRetain bool
		var sslCert, sslKey, sslCA, sslServer, willTopic, willMsg string
		for _, v := range opts {
			kv := strings.Split(v, "=")
			if len(kv) != 2 {
				println(v, "option should be key=value")
				return nil, err
			}

			switch kv[0] {
			case "clean":
				options = append(options, mqtt.WithCleanSession(kv[0] == "y"))
			case "ssl":
				ssl = kv[1] == "y"
			case "ssl_cert":
				sslCert = kv[1]
			case "ssl_key":
				sslKey = kv[1]
			case "ssl_ca":
				sslCA = kv[1]
			case "ssl_server":
				sslServer = kv[1]
			case "ssl_skip_verify":
				sslSkipVerify = kv[1] == "y"
			case "will":
				will = kv[1] == "y"
			case "will_topic":
				willTopic = kv[1]
			case "will_qos":
				qos, err := strconv.Atoi(kv[1])
				if err != nil || qos > 2 {
					utils.InvalidQos()
					return nil, err
				}
				willQos = mqtt.QosLevel(qos)
			case "will_msg":
				willMsg = kv[1]
			case "will_retain":
				willRetain = kv[1] == "y"
			}
		}
		if ssl {
			options = append(options, mqtt.WithTLS(sslCert, sslKey, sslCA, sslServer, sslSkipVerify))
		}
		if will {
			options = append(options, mqtt.WithWill(willTopic, willQos, willRetain, []byte(willMsg)))
		}
	}
	options = append(options, mqtt.WithClientID(args[1]))
	options = append(options, mqtt.WithIdentity(args[2], args[3]))
	keepAlive, err := strconv.Atoi(args[4])
	if err != nil {
		return nil, err
	}
	options = append(options, mqtt.WithKeepalive(uint16(keepAlive), 1.2))
	options = append(options, mqtt.WithVersion(version, false))
	return newClient(options, args[0], info)
}

func execDisConn(args []string) bool {
	var client mqtt.Client
	if client != nil {
		client.Destroy(!(len(args) > 0 && args[1] != "force"))
	}

	os.Exit(0)

	return true
}

type ConnectionPreProcess func(client mqtt.Client)
type PostHandler func(client mqtt.Client)

func newClient(options []mqtt.Option, server string, info *domain.ClientInfo) (client mqtt.Client, err error) {

	allOpts := append([]mqtt.Option{
		mqtt.WithPubHandleFunc(utils.PubHandler),
		mqtt.WithConnHandleFunc(utils.ConnHandler),
		mqtt.WithUnsubHandleFunc(utils.UnSubHandler),
		mqtt.WithNetHandleFunc(utils.NetHandler),
		mqtt.WithSubHandleFunc(utils.SubHandler),
		mqtt.WithRouter(mqtt.NewRegexRouter()),
	}, options...)

	client, err = mqtt.NewClient(allOpts...)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	// add client info
	domain.ClientInfoMap[client] = info
	utils.AddHandler(utils.WelcomeTopic, welcomeHandler)
	utils.AddHandler(utils.CmdTopic, cmdHandler)
	client.HandleTopic(".*", func(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
		if clientInfo, ok := domain.ClientInfoMap[client]; ok {
			switch clientInfo.Model {
			case domain.SM4:
				var token []byte
				if utils.WelcomeTopic == topic {
					token = []byte(clientInfo.Token)
				} else {
					token = []byte(clientInfo.Welcome().CryptoSecret)
				}
				decrypt, err := crypto.DoSM4Decrypt(msg, token, token)
				if err != nil {
					fmt.Printf("msg DoSM4Decrypt failed ,close client,msg=%s \n", base64.StdEncoding.EncodeToString(msg))
					fmt.Printf(domain.LineStart)
					return
				}
				fmt.Printf("\n[%v] message: %v qos:%v \n", topic, string(decrypt), qos)
				fmt.Printf(domain.LineStart)
				utils.HandleMsg(client, topic, qos, decrypt)
			default:
				fmt.Printf("\n[%v] message: %v qos:%v \n", topic, string(msg), qos)
				fmt.Printf(domain.LineStart)
				utils.HandleMsg(client, topic, qos, msg)
			}
		} else {
			fmt.Printf("clientInfo not exist ,close client \n")
			fmt.Printf(domain.LineStart)
		}

	})

	//client.HandleTopic("sys/welcome", welcomeHandler)
	//client.HandleTopic("sys/cmd/*", cmdHandler)
	err = client.ConnectServer(server, allOpts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

var welcomeHandler = func(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
	if clientInfo, ok := domain.ClientInfoMap[client]; ok {
		saveWelcome(msg, clientInfo)
	} else {

		fmt.Printf("clientInfo not exist ,close client")
		fmt.Printf(domain.LineStart)
	}

}
var cmdHandler = func(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
	split := strings.Split(topic, "/")
	cmdTag := split[len(split)-1]
	completeTopic := strings.Join([]string{"sys", "cmdSync", cmdTag, "complete"}, "/")
	pubMsg := CreateSinglePubMsg(0, completeTopic, "")
	client.Publish(pubMsg...)
}

func saveWelcome(msg []byte, clientInfo *domain.ClientInfo) {
	var welcome domain.WelcomeInfo
	err1 := json.Unmarshal(msg, &welcome)
	if err1 != nil {
		fmt.Printf("welcome msg Unmarshal failed ")
		fmt.Printf(domain.LineStart)

	}
	clientInfo.SetWelcome(welcome)

}

func connUsage() {
	fmt.Println(`c, conn [server:port] sn gk pwd keepalive - connect to server`)
}
