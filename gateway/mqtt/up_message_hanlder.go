package mqtt

import (
	"fmt"
	"github.com/goiiot/libmqtt"
)

func PubHandler(client libmqtt.Client, topic string, msg string, err error) {

	if err != nil {
		m := fmt.Sprintf("pub %v %v failed,error=%v", topic, msg, err)
		fmt.Println(m)
	} else {
		m := fmt.Sprintf("pub %v %v success", topic, msg)
		fmt.Println(m)
	}
}

func ConnHandler(client libmqtt.Client, server string, code byte, err error) {
	if err != nil {
		fmt.Printf("\nconnect to server:%v error:%v\n", server, err)
	} else if code != libmqtt.CodeSuccess {
		fmt.Printf("\nconnection rejected by server：%v, code:%v\n", server, code)
	} else {
		fmt.Printf("\nconnected to server:%v\n", server)
	}
}
func UnSubHandler(client libmqtt.Client, topics []string, err error) {
	if err != nil {
		fmt.Println("\nunsub", topics, "failed, error =", err)
	} else {
		fmt.Println("\nunsub", topics, "success")
	}
}
func SubHandler(client libmqtt.Client, topics []*libmqtt.Topic, err error) {
	if err != nil {
		fmt.Println("\nsub", topics, "failed, error =", err)
	} else {
		fmt.Println("\nsub", topics, "success")
	}
}
func NetHandler(client libmqtt.Client, server string, err error) {
	fmt.Println("\nconnection to server, error:", err)
	client.Destroy(false)
}

func NetHandlerWithReconnect(client libmqtt.Client, server string, err error) {
	fmt.Println("\nconnection to server, error:", err)
	client.Reconnect(server)
}
