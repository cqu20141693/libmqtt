package mqtt

import (
	"fmt"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/gateway/initialize/logger/cclog"
	"time"
)

func NewClient(options []mqtt.Option, server string) (client mqtt.Client, err error) {

	allOpts := append([]mqtt.Option{
		// 处理up message 异常情况
		mqtt.WithPubHandleFunc(PubHandler),
		mqtt.WithReceiveHandleFunc(ReceiveHandler),
		mqtt.WithConnHandleFunc(ConnHandler),
		mqtt.WithUnsubHandleFunc(UnSubHandler),
		mqtt.WithNetHandleFunc(NetHandlerWithReconnect),
		mqtt.WithSubHandleFunc(SubHandler),
		mqtt.WithRouter(mqtt.NewRegexRouter()),
		// 支持连接失败，自动重连
		mqtt.WithAutoReconnect(true),
		mqtt.WithBackoffStrategy(1*time.Second, 10*time.Second, 1.5),
	}, options...)
	client, err = mqtt.NewClient(allOpts...)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	utils.AddHandler(utils.WelcomeTopic, welcomeHandler)
	utils.AddHandler(utils.CmdTopic, cmdHandler)
	client.HandleTopic(".*", func(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
		// 处理原始消息
		// 如果存在加解密可以先处理
		fmt.Printf("\n[%v] message: %v qos:%v \n", topic, string(msg), qos)
		utils.HandleMsg(client, topic, qos, msg)

	})
	err = client.ConnectServer(server, allOpts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func cmdHandler(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {

}

var welcomeHandler = func(client mqtt.Client, topic string, qos mqtt.QosLevel, msg []byte) {
	cclog.Info("receive welcome msg topic=%s qos=%d msg=%v", topic, qos, string(msg))

}
