package simulator

import (
	"encoding/json"
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/random"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"strconv"
	"strings"
)

type SimulateClient struct {
	ConnectionConfig map[string]interface{}
}

// Close 关闭客户端连接和任务
//
//	@receiver c
//	@param channelId
func (c *SimulateClient) Close(channelId string) {
	cclog.Info(fmt.Sprintf("simulator client:%v already closed", channelId))
}

// Open 打开通道连接
//
//	@receiver c
//	@param channelId
func (c *SimulateClient) Open(channelId string) {
	cclog.Info(fmt.Sprintf("simulator client:%v already started", channelId))
}

// ReadTagValues  读设备取点位值
//
//	@receiver c
//	@param device
//	@param tags
//	@return connectors.ReadTagResult
func (c *SimulateClient) ReadTagValues(device connectors.Device, tags []connectors.Tag) connectors.ReadTagResult {

	var result = connectors.ReadTagResult{
		DeviceId: device.DeviceID,
		Fails:    map[string]string{},
		Success:  map[string]interface{}{},
	}
	for _, tag := range tags {
		success, v := c.GetValue(tag)
		if success {
			result.Success[tag.TagId] = v
		} else {
			result.Fails[tag.TagId] = v.(string)
		}
	}
	return result
}
func (c *SimulateClient) GetValue(tag connectors.Tag) (success bool, value interface{}) {
	addressType := tag.Extension["address"].(string)
	dataType := tag.DataType
	if addressType == "Static" {
		if tag.Value != nil {
			return true, tag.Value
		} else {
			defaultV := tag.Extension["default"]
			if defaultV == nil {
				tag.Value = getRandom(dataType)
				return true, tag.Value
			}
			switch dataType {
			case constants.STRING:
				tag.Value = defaultV.(string)
				return true, tag.Value
			case constants.INTEGER:
				switch defaultV.(type) {
				case float64:
					tag.Value = int64(defaultV.(float64))
					return true, tag.Value
				default:
					tag.Value, _ = strconv.Atoi(defaultV.(string))
					return true, tag.Value
				}
			case constants.LONG:
				switch defaultV.(type) {
				case float64:
					tag.Value = int64(defaultV.(float64))
					return true, tag.Value
				default:
					tag.Value, _ = strconv.ParseInt(defaultV.(string), 10, 64)
					return true, tag.Value
				}
			case constants.BOOL:
				switch defaultV.(type) {
				case bool:
					tag.Value = defaultV.(bool)
					return true, tag.Value
				default:
					tag.Value, _ = strconv.ParseBool(defaultV.(string))
					return true, tag.Value
				}

			case constants.Float:
				switch defaultV.(type) {
				case float64:
					tag.Value = defaultV.(float64)
					return true, tag.Value
				default:
					tag.Value, _ = strconv.ParseFloat(defaultV.(string), 64)
					return true, tag.Value
				}
			default:
				return false, fmt.Sprintf("not support dataType:%v", dataType)
			}
		}
	} else if addressType == "Random" {
		v := getRandomWithConfig(dataType, tag.Extension)
		if v == nil {
			return false, fmt.Sprintf("not support dataType:%v", dataType)
		}
		return true, v
	} else if addressType == "Increment" {
		min := tag.Extension["min"].(float64)
		max := tag.Extension["max"].(float64)
		step := tag.Extension["step"].(float64)
		if tag.Value == nil {
			tag.Value = min
			return true, min
		} else {
			v := min + step
			if v > max {
				tag.Value = min - step
				return true, max
			}
			if v < min {
				tag.Value = min - step
				return true, min
			}
			return true, v
		}

	} else if addressType == "Enum" {
		enums := tag.Extension["enums"].(string)
		values := strings.Split(enums, ",")
		return true, random.Choice(values)
	} else if addressType == "JsonArray" {
		jsonData := tag.Extension["jsonArray"].(string)
		// 定义一个切片来保存未知类型的JSON数据
		var unknownData []json.RawMessage
		// 解析JSON数组
		err := json.Unmarshal([]byte(jsonData), &unknownData)
		if err != nil {
			fmt.Println("解析JSON数组失败:", err)
			return false, "解析JSON数组失败"
		}
		// 遍历未知类型的JSON数据
		var jsonArray []interface{}
		for _, rawMsg := range unknownData {
			var item interface{}

			// 解析每个元素
			err := json.Unmarshal(rawMsg, &item)
			if err != nil {
				fmt.Println("解析JSON元素失败:", err)
				continue
			}
			jsonArray = append(jsonArray, item)
		}
		return true, random.Choice(jsonArray)
	}
	return false, fmt.Sprintf("not support simulator address type:%v", addressType)
}
func getRandom(dataType constants.DataType) interface{} {
	switch dataType {
	case constants.STRING:
		return random.RandString(8)
	case constants.INTEGER:
		return random.Int31()
	case constants.LONG:
		return random.Int63()
	case constants.BOOL:
		return random.Choice([]bool{true, false})
	case constants.Float:
		return random.Float64()
	default:
		return random.RandString(8)
	}
}
func getRandomWithConfig(dataType constants.DataType, extension map[string]interface{}) interface{} {
	switch dataType {
	case constants.STRING:
		return random.RandString(8)
	case constants.INTEGER:
		min := utils.GetIntFromViper(extension, "min")
		max := utils.GetIntFromViper(extension, "max")
		return random.RandInt(min, max)
	case constants.LONG:
		min := utils.GetInt64FromViper(extension, "min")
		max := utils.GetInt64FromViper(extension, "max")
		return random.RandInt64(min, max)
	case constants.BOOL:
		return random.Choice([]bool{true, false})
	case constants.Float:
		min := extension["min"].(float64)
		max := extension["max"].(float64)
		return random.RandFloat64(min, max)
	default:
		fmt.Println(fmt.Sprintf("not support dataType:%v", dataType))
		return nil
	}
}
