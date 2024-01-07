package mqtt

import (
	"fmt"
	"github.com/goiiot/libmqtt"
)

func ReceiveHandler(client libmqtt.Client, topic string, msg string, err error) {

	if err != nil {
		m := fmt.Sprintf("receive %v %v failed,error=%v", topic, msg, err)
		fmt.Println(m)
	} else {
		m := fmt.Sprintf("receive %v %v success", topic, msg)
		fmt.Println(m)
	}
}
