package mqtt

import (
	"fmt"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
)

func PubHandler(client libmqtt.Client, topic string, msg string, err error) {

	if err != nil {
		m := fmt.Sprintf("client=%s pub %v %v failed,error=%v", client.ClientId(), topic, msg, err)
		cclog.SugarLogger.Error(m)
	} else {
		if client.PubMetric != nil {
			client.PubMetric.Inc(1)
		}
	}
}

func UnSubHandler(client libmqtt.Client, topics []string, err error) {
	if err != nil {
		cclog.SugarLogger.Info(fmt.Sprintf("client=%s unsub %v failed, error:%v", client.ClientId(), topics, err))
	} else {
		cclog.SugarLogger.Info(fmt.Sprintf("client=%s unsub %v success", client.ClientId(), topics))
	}
}
func SubHandler(client libmqtt.Client, topics []*libmqtt.Topic, err error) {
	if err != nil {
		cclog.SugarLogger.Info(fmt.Sprintf("client=%s sub %v failed, error:%v", client.ClientId(), topics, err))
	} else {
		cclog.SugarLogger.Info(fmt.Sprintf("client=%s sub %v success", client.ClientId(), topics))
	}
}
