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
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/crypto"
	"github.com/goiiot/libmqtt/cmd/domain"
	"github.com/goiiot/libmqtt/cmd/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	flag.Parse()
	osCh := make(chan os.Signal, 2)
	signal.Notify(osCh, os.Kill, os.Interrupt)
	wg := &sync.WaitGroup{}

	// handle system signals
	go func() {
		for range osCh {
			os.Exit(0)
		}
	}()
	// handle user input
	wg.Add(1)
	go func() {
		fmt.Printf(domain.LineStart)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" {
				args := strings.Split(text, " ")
				if text == "quit" {
					execCmd(args)
					break
				}
				execCmd(args)
			}
			fmt.Printf(domain.LineStart)
		}
		wg.Done()
	}()
	wg.Wait()
}

func execCmd(params []string) {
	cmd := strings.ToLower(params[0])
	args := make([]string, len(params)-1)
	copy(args, params[1:])
	ok := false
	switch cmd {
	// support 网关和普通设备连接，支持app订阅连接
	case "c", "conn":
		ok = handleConn(params, args, mqtt.V311)
	case "m-c", "m-conn":
		infos, err := getMirrorConnInfo(args)
		if err == nil {
			ok = handleConn(params, infos, mqtt.V311)
		} else {
			fmt.Println(`m-conn,m-c server getMirrorInfoAPI keepalive  - Mirror login`)
		}
	case "s-conn":
		signature := crypto.DoSignature(params[4])
		if signature == "" {
			fmt.Println(`s-conn signature failed  - signature login`)
			break
		}
		args[3] = signature
		ok = handleConn(params, args, mqtt.V311)
	case "sg-c", "sg-conn":
		infos, err := getSubGroupConnInfo(args)
		if err == nil {
			ok = handleConn(params, infos, mqtt.V311)
		} else {
			fmt.Println(`sg-conn,sg-c server getSubGroupInfoAPI [device1,device2] keepalive  - subGroup login`)
		}
		// support 网关和普通设备连接，支持app订阅连接
	case "conn5":
		ok = handleConn(params, args, mqtt.V5)
	case "p", "pub":
		if client, exit := domain.ClientMap[args[len(args)-1]]; exit {
			ok = execPub(client, args[:len(args)-1])
		} else {
			print("clientId not exist,please use the lookup command to view the client")
		}
	case "s", "sub":
		if client, exit := domain.ClientMap[args[len(args)-1]]; exit {
			ok = execSub(client, args[:len(args)-1])
		} else {
			print("clientId not exist,please use the lookup command to view the client")
		}
	case "un", "unsub":
		if len(args) == 2 {
			if client, exit := domain.ClientMap[args[len(args)-1]]; exit {
				ok = execUnSub(client, args[:len(args)-1])
			} else {
				print("clientId not exist,please use the lookup command to view the client")
			}
		}
	case "q", "exit":
		for _, client := range domain.ClientMap {
			client.Destroy(false)
			delete(domain.ClientMap, args[0])
			delete(domain.CmdMap, args[0])
			ok = true
		}
	case "d", "disconnect":
		// d 1
		// d 1 force
		if len(args) >= 1 {
			if client, exit := domain.ClientMap[args[0]]; exit {
				force := len(args) == 2 && args[1] == "force"
				client.Destroy(force)
				delete(domain.ClientMap, args[0])
				delete(domain.CmdMap, args[0])
				ok = true
			} else {
				print("clientId not exist,please use the lookup command to view the client")
			}
		}
	case "l", "lookup":
		ok = lookup()
	//ok = execDisConn(args)
	case "db", "decrypt-base64":
		// db SM4 password base64-data
		if len(args) == 3 {
			switch args[0] {
			case domain.SM4:
				cipherText, err := base64.StdEncoding.DecodeString(args[2])
				if err != nil {
					fmt.Println("base64 decode failed")
					break
				}
				decrypt, err := crypto.DoSM4Decrypt(cipherText, []byte(args[1]), []byte(args[1]))
				if err != nil {
					fmt.Println("Do SM4 Decrypt failed")
					break
				}
				fmt.Printf("decrypt data=%s  \n", string(decrypt))
				ok = true
			}
		}
	case "cb", "crypto-base64":
		// cb SM4 password json data
		if len(args) == 4 {
			switch args[0] {
			case domain.SM4:
				bytes, err := getBytes(args[2], args[3])
				if err != nil {
					fmt.Println("data encode type error")
					return
				}
				sm4, err := crypto.DoSM4(bytes, []byte(args[1]), []byte(args[1]))
				if err != nil {
					fmt.Println("Do SM4 Crypt failed")
					return
				}
				fmt.Printf("base64 sm4 crypt data=%s\n", base64.StdEncoding.EncodeToString(sm4))
				ok = true
			}
		}
	case "ch", "crypto-hex":
		//ch SM4 password json data
		if len(args) == 3 {
			switch args[0] {
			case domain.SM4:
				sm4, err := crypto.DoSM4([]byte(args[2]), []byte(args[1]), []byte(args[1]))
				if err != nil {
					fmt.Println("Do SM4 Crypt failed")
					return
				}
				fmt.Printf("base64 sm4 crypt data=%s\n", hex.EncodeToString(sm4))
			}
		}
	case "st", "signature-token":
		//st SM3 token
		if len(args) == 2 {
			fmt.Printf("sm3 signature data=%s\n", crypto.DoSignature(strings.Join(args, ":")))
			ok = true
		}

	case "?", "h", "help":
		ok = usage()
	}

	if !ok {
		usage()
	}
}

func getCipherData(client mqtt.Client, encodeType, data string) (r []byte, err error) {
	bytes, err := getBytes(encodeType, data)
	if err != nil {
		return nil, err
	}
	if info, exit := domain.ClientInfoMap[client]; exit {
		switch info.Model {
		case domain.SM4:
			key := []byte(info.Welcome().CryptoSecret)
			return crypto.DoSM4(bytes, key, key)
		case domain.AES:
		}
	}
	return bytes, err
}

func handleConn(params []string, args []string, version mqtt.ProtoVersion) (ok bool) {
	split := strings.Split(params[4], ":")
	var model, token string
	if len(split) > 1 {
		model = split[0]
		token = split[1]
	} else {
		token = split[0]
	}
	info := domain.NewClientInfo(params[0], params[1], params[2], params[3], model, token, "", params[4:])
	// exec conn
	client, err := execConn(args, version, info)
	if err == nil {
		counter := atomic.AddInt64(&domain.IdGenerator, 1)
		clientId := strconv.FormatInt(counter, 10)
		info.SetId(clientId)
		domain.ClientMap[clientId] = client
		domain.CmdMap[clientId] = info
		ok = true
	} else {
		delete(domain.ClientInfoMap, client)
	}
	return
}

func getSubGroupConnInfo(args []string) (infos []string, err error) {
	if len(args) != 6 {
		return nil, err
	}
	devices := strings.Split(args[4], ",")
	url := args[3] + "?" + "groupKey=" + args[1] + "&" + "appKey=" + args[2]
	mqttInfo, err := http.GetSubMqttInfo(url, devices)
	if err != nil {
		return nil, err
	}
	// server,clientId,username,pwd,keepalive
	return []string{args[0], mqttInfo.ClientIdentifier, mqttInfo.Username, "SG:" + mqttInfo.Password, args[5]}, nil
}
func getMirrorConnInfo(args []string) (infos []string, err error) {
	if len(args) != 5 {
		return nil, errors.New("params length error")
	}
	mqttInfo, err := http.GetMirrorMqttInfo(args[1:4])
	if err != nil {
		return nil, err
	}
	// server,clientId,username,pwd,keepalive
	return []string{args[0], mqttInfo.ClientIdentifier, mqttInfo.Username, "M:" + mqttInfo.Password, args[4]}, nil
}

func lookup() bool {
	for s, cmd := range domain.CmdMap {
		fmt.Printf("clientId=%v cmd=%v \n", s, cmd)
	}
	return true
}

func usage() bool {
	fmt.Println("Usage:")

	connUsage()
	pubUsage()
	subUsage()
	unSubUsage()
	fmt.Println("d clientId [force] - disconnect client")
	fmt.Println(`q, exit - exit CLI`)
	fmt.Println(`h, help - print this help message`)
	return true
}
