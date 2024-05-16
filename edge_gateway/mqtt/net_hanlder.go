package mqtt

import (
	"fmt"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
)

func ConnHandler(client libmqtt.Client, server string, code byte, err error) {
	if err != nil {
		cclog.SugarLogger.Error(fmt.Sprintf("client=%s connect to server:%v error:%v", client.ClientId(), server, err))
	} else if code != libmqtt.CodeSuccess {
		cclog.SugarLogger.Error(fmt.Sprintf("client=%s connection rejected by server：%v, code:%v", client.ClientId(), server, code))
	} else {
		cclog.SugarLogger.Info(fmt.Sprintf("client=%s connected to server:%v", client.ClientId(), server))
	}
}

func NetHandler(client libmqtt.Client, server string, err error) {
	cclog.SugarLogger.Info(fmt.Sprintf("client=%s connection to server, error:%v", client.ClientId(), err))
	client.Destroy(false)
}

func NetHandlerWithReconnect(client libmqtt.Client, server string, err error) {

	cclog.SugarLogger.Info(fmt.Sprintf("client=%s will reconnection to server, error:%v", client.ClientId(), err))
	client.Reconnect(server)
}
