package utils

import mqtt "github.com/goiiot/libmqtt"

func CreatePublishPacket(topic string, qos byte, message string) *mqtt.PublishPacket {
	return &mqtt.PublishPacket{
		TopicName: topic,
		Qos:       qos,
		Payload:   []byte(message),
	}
}
