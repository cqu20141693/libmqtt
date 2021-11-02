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
	"flag"
	"fmt"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	lineStart = "> "
)

var idGenerator int64
var clientMap = make(map[string]*mqtt.AsyncClient, 8)
var cmdMap = make(map[string][]string, 8)

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
		fmt.Printf(lineStart)
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
			fmt.Printf(lineStart)
		}
		wg.Done()
	}()
	wg.Wait()
}

func execCmd(params []string) {
	cmd := strings.ToLower(params[0])
	args := params[1:]
	ok := false
	switch cmd {
	// support 网关和普通设备连接，支持app订阅连接
	case "c", "conn":
		ok = handleConn(params, args, ok)
	case "m-c", "m-conn":
		infos, err := getMirrorConnInfo(args)
		if err == nil {
			ok = handleConn(params, infos, ok)
		} else {
			fmt.Println(`m-conn,m-c server getMirrorInfoAPI keepalive  - Mirror login`)
		}
	case "sg-c", "sg-conn":
		infos, err := getSubGroupConnInfo(args)
		if err == nil {
			ok = handleConn(params, infos, ok)
		} else {
			fmt.Println(`sg-conn,sg-c server getSubGroupInfoAPI [device1,device2] keepalive  - subGroup login`)
		}
	case "p", "pub":
		if client, exit := clientMap[args[len(args)-1]]; exit {
			ok = execPub(client, args[:len(args)-1])
		} else {
			print("clientId not exist,please use the lookup command to view the client")
		}
	case "s", "sub":
		if client, exit := clientMap[args[len(args)-1]]; exit {
			ok = execSub(client, args[:len(args)-1])
		} else {
			print("clientId not exist,please use the lookup command to view the client")
		}
	case "u", "unsub":
		if client, exit := clientMap[args[len(args)-1]]; exit {
			ok = execUnSub(client, args[:len(args)-1])
		} else {
			print("clientId not exist,please use the lookup command to view the client")
		}
	case "q", "exit":
		for _, client := range clientMap {
			client.Destroy(false)
			delete(clientMap, args[0])
			delete(cmdMap, args[0])
			ok = true
		}
	case "d", "disconnect":
		if client, exit := clientMap[args[0]]; exit {
			force := len(args) == 2 && args[1] == "force"
			client.Destroy(force)
			delete(clientMap, args[0])
			delete(cmdMap, args[0])
			ok = true
		} else {
			print("clientId not exist,please use the lookup command to view the client")
		}
	case "l", "lookup":
		ok = lookup()
	//ok = execDisConn(args)
	case "?", "h", "help":
		ok = usage()
	}

	if !ok {
		usage()
	}
}

func handleConn(params []string, args []string, ok bool) bool {
	client, err := execConn(args)
	if err == nil {
		counter := atomic.AddInt64(&idGenerator, 1)
		clientId := strconv.FormatInt(counter, 10)
		clientMap[clientId] = client
		cmdMap[clientId] = params
		ok = true
	}
	return ok
}

func getSubGroupConnInfo(args []string) (infos []string, err error) {
	if len(args) != 4 {
		return nil, err
	}
	devices := strings.Split(args[2], ",")
	mqttInfo, err := http.GetSubMqttInfo(args[1], devices)
	if err != nil {
		return nil, err
	}
	// server,clientId,username,pwd,keepalive
	return []string{args[0], mqttInfo.ClientIdentifier, mqttInfo.Username, "SG:" + mqttInfo.Password, args[3]}, nil
}
func getMirrorConnInfo(args []string) (infos []string, err error) {
	if len(args) != 3 {
		return nil, err
	}
	mqttInfo, err := http.GetMirrorMqttInfo(args[1])
	if err != nil {
		return nil, err
	}
	// server,clientId,username,pwd,keepalive
	return []string{args[0], mqttInfo.ClientIdentifier, mqttInfo.Username, "M:" + mqttInfo.Password, args[2]}, nil
}

func lookup() bool {
	for s, cmd := range cmdMap {
		fmt.Printf("clientId=%v cmd=%v \n", s, cmd)
	}
	return true
}

func usage() bool {
	print("Usage\n\n")
	print("  ")
	connUsage()
	print("  ")
	pubUsage()
	print("  ")
	subUsage()
	print("  ")
	unSubUsage()
	println(`  q, exit [force] - disconnect and exit`)
	println(`  h, help - print this help message`)
	return true
}
