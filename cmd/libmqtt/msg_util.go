package main

import (
	"encoding/binary"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/cmd/domain"
	"math"
	"strconv"
)

func CreateSinglePubMsg(qos byte, topicName, payload string) []*mqtt.PublishPacket {
	pubs := make([]*mqtt.PublishPacket, 1)
	pubs = append(pubs, &mqtt.PublishPacket{
		TopicName: topicName,
		Qos:       mqtt.QosLevel(qos),
		Payload:   []byte(payload),
	})
	return pubs
}

/*
	根据编码类型序列化字节数组
*/
func getBytes(encodeType, dataStr string) (bytes []byte, err error) {
	switch encodeType {
	case domain.String, domain.Json, domain.Bin:
		return []byte(dataStr), nil
	case domain.Int:
		parseInt, err := strconv.ParseInt(dataStr, 10, 32)
		if err != nil {
			return nil, err
		}
		return []byte{byte(parseInt >> 24), byte(parseInt >> 16), byte(parseInt >> 8), byte(parseInt)}, nil
	case domain.Long:
		long, err := strconv.ParseInt(dataStr, 10, 64)
		if err != nil {
			return nil, err
		}
		return []byte{byte(long >> 56), byte(long >> 48), byte(long >> 40), byte(long >> 32), byte(long >> 24), byte(long >> 16), byte(long >> 8), byte(long)}, nil
	case domain.Float:
		float, err := strconv.ParseFloat(dataStr, 32)
		if err != nil {
			return nil, err
		}
		bits := math.Float32bits(float32(float))
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, bits)
		return bytes, nil
	case domain.Double:
		float, err := strconv.ParseFloat(dataStr, 64)
		if err != nil {
			return nil, err
		}
		bits := math.Float64bits(float)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, bits)
		return bytes, nil

	}
	return nil, err
}
