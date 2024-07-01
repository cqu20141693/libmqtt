package mqtt

import (
	"fmt"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"time"
)

func init() {
	// 初始化消息处理
	utils.AddHandler(utils.WelcomeTopic, welcomeHandler)
	utils.AddHandler(utils.CmdTopic, cmdHandler)

}
func CreatClient(info *domain.MqttClientAddInfo) (libmqtt.Client, error) {
	options := make([]libmqtt.Option, 0)
	options = append(options, libmqtt.WithCleanSession(true))
	options = append(options, libmqtt.WithClientID(info.ClientID))
	options = append(options, libmqtt.WithIdentity(info.Username, info.Password))
	options = append(options, libmqtt.WithKeepalive(uint16(info.Keepalive), 1.2))
	options = append(options, libmqtt.WithVersion(libmqtt.V311, false))
	// 支持日志等级
	options = append(options, libmqtt.WithLog(info.LogLevel))

	client, err := NewClient(options, info.Address)
	return client, err
}

func NewClient(options []libmqtt.Option, server string) (client libmqtt.Client, err error) {

	allOpts := append([]libmqtt.Option{
		// 处理up message 异常情况
		libmqtt.WithLog(libmqtt.Info),
		//
		libmqtt.WithPubHandleFunc(PubHandler),
		libmqtt.WithReceiveHandleFunc(ReceiveHandler),
		libmqtt.WithConnHandleFunc(ConnHandler),
		libmqtt.WithUnsubHandleFunc(UnSubHandler),
		libmqtt.WithNetHandleFunc(NetHandlerWithReconnect),
		libmqtt.WithSubHandleFunc(SubHandler),
		libmqtt.WithRouter(libmqtt.NewRegexRouter()),
		// 支持连接失败，自动重连
		libmqtt.WithAutoReconnect(true),
		libmqtt.WithBackoffStrategy(1*time.Second, 10*time.Second, 1.5),
	}, options...)
	client, err = libmqtt.NewClient(allOpts...)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	client.HandleTopic(".*", func(client libmqtt.Client, topic string, qos libmqtt.QosLevel, msg []byte) {
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

func cmdHandler(client libmqtt.Client, topic string, qos libmqtt.QosLevel, msg []byte) {
	cclog.Info("receive cmd msg client=%s topic=%s qos=%d msg=%v", client.ClientId(), topic, qos, string(msg))
}

var welcomeHandler = func(client libmqtt.Client, topic string, qos libmqtt.QosLevel, msg []byte) {
	cclog.Info("receive welcome msg topic=%s qos=%d msg=%v", topic, qos, string(msg))

}
