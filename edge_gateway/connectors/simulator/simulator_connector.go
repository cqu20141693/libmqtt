package simulator

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"time"
)

type SimulateConnector struct {
	connectors.ChannelInstance
	Config map[string]interface{} // 整体配置
	// 连接配置，根据协议不同，连接配置不同
	ConnectConfig map[string]interface{}
	Client        SimulateClient
}

func (s *SimulateConnector) ServerSideRpcHandlerSync(eventData connectors.RpcEvent) connectors.RpcResult {
	//TODO implement me
	panic("implement me")
}
func (s *SimulateConnector) ServerSideRpcHandler(event connectors.RpcEvent) {
	//TODO implement me
	panic("implement me")
}

func (s *SimulateConnector) AttributesUpdateSync(eventData connectors.UpdateTagEvent) connectors.UpdateTagResult {
	//TODO implement me
	panic("implement me")
}
func (s *SimulateConnector) AttributesUpdate(event connectors.UpdateTagEvent) {
	//TODO implement me
	panic("implement me")
}

func (s *SimulateConnector) AttributesReadSync(eventData connectors.ReadTagEvent) connectors.ReadTagResult {
	var ret = connectors.ReadTagResult{
		DeviceId: eventData.DeviceId,
		Fails:    map[string]string{},
		Success:  map[string]interface{}{},
	}
	var device connectors.Device
	if _, exist := s.ChannelCache.Devices[eventData.DeviceId]; !exist {
		for _, tagId := range eventData.TagIds {
			ret.Fails[tagId] = "设备不存在"
		}
		return ret
	}

	//1. 判断通道连接
	if !s.IsConnected() {
		for _, tagId := range eventData.TagIds {
			ret.Fails[tagId] = "通道未连接"
		}
		return ret
	}

	// 2. 校验属性可读性和存在性
	tags := make([]connectors.Tag, 0)
	for _, tagId := range eventData.TagIds {
		tag := device.FindTag(tagId)
		if tag.TagId != tagId {
			ret.Fails[tagId] = "点位不存在"
		}
		if tag.TagType != constants.ATTRIBUTES {
			ret.Fails[tagId] = "遥测点位不可读"
		}
		if tag.ReadWrite&constants.ReadOnly != 1 {
			ret.Fails[tagId] = "属性值不可读"
		}
		tags = append(tags, tag)
	}
	if len(tags) > 0 {
		tagValues := s.Client.ReadTagValues(device, tags)
		for tagId, reason := range tagValues.Fails {
			ret.Fails[tagId] = reason
		}
		for tagId, value := range tagValues.Success {
			ret.Success[tagId] = value
		}
	}
	return ret

}

func (s *SimulateConnector) AttributesRead(eventData connectors.ReadTagEvent) {
	result := s.AttributesReadSync(eventData)
	event := events.SouthEvent{
		Type:       events.Telemetry,
		EventTopic: events.TelemetryTopic,
		DeviceId:   result.DeviceId,
		Data:       result,
	}
	s.ConnectorManager.SouthEventCh <- event
}

func (s *SimulateConnector) Close(remote bool) {
	s.Client.Close(s.ChannelId)
	s.SetConnected(false, false)
	s.Scheduler.Stop()

}

func (s *SimulateConnector) IsConnected() bool {
	return s.Connected
}

func (s *SimulateConnector) GetId() string {

	return s.ChannelId
}
func (s *SimulateConnector) LoadConfig() {
	// 定义设备的模型
	devices := make(map[string]connectors.Device)
	tagMap := make(map[string]connectors.Tag)
	var tagCount uint16
	deviceConfig := s.Config["devices"].([]interface{})
	for _, conf := range deviceConfig {
		config := conf.(map[string]interface{})
		deviceID := config["deviceId"].(string)
		deviceName := config["deviceName"].(string)
		deviceType := config["deviceType"].(string)

		sendDataOnlyOnChange := config["sendDataOnlyOnChange"].(bool)
		var extension = make(map[string]interface{})
		if e, exist := config["extension"]; exist {
			extension = e.(map[string]interface{})
		}
		// 解析Tag
		tagConfig := make([]interface{}, 0)
		if ts, exist := config["tags"]; exist {
			tagConfig = ts.([]interface{})
		}
		tags := make([]connectors.Tag, 0, len(tagConfig))
		for _, tagConf := range tagConfig {
			tf := tagConf.(map[string]interface{})
			tagCount += 1
			name := tf["name"].(string)
			tagId := tf["tagId"].(string)
			dataType := constants.DataType(tf["dataType"].(string))
			tagType := constants.TagType(tf["tagType"].(string))
			frequency := int64(tf["frequency"].(float64))
			tagSendDataOnlyOnChange := tf["sendDataOnlyOnChange"].(bool)
			tagExtension := make(map[string]interface{})
			if ext, exist := tf["extension"]; exist {
				tagExtension = ext.(map[string]interface{})
			}

			var readWrite constants.ReadWriteType
			if rw, exist := tf["readWrite"]; exist {
				readWrite = constants.ReadWriteType(rw.(float64))
			}
			tag := connectors.NewTag(name, tagId, deviceID, dataType, readWrite, tagType, frequency, tagSendDataOnlyOnChange, tagExtension)
			// 静态值
			if v, exist := tf["value"]; exist {
				tag.Value = v
			}
			tagMap[tagId] = tag
			tags = append(tags, tag)
		}
		functionConfig := make([]interface{}, 0)
		if fs, exist := config["functions"]; exist {
			functionConfig = fs.([]interface{})
		}

		// 解析功能
		functions := make([]connectors.Function, 0, len(functionConfig))
		for _, functionConf := range functionConfig {
			f := functionConf.(map[string]interface{})
			fe := make(map[string]interface{})
			if ext, exist := f["extension"]; exist {
				fe = ext.(map[string]interface{})
			}

			function := connectors.Function{
				Method:    f["method"].(string),
				Name:      f["name"].(string),
				DeviceID:  deviceID,
				Extension: fe,
			}
			functions = append(functions, function)
		}

		device := connectors.NewDevice(deviceID, deviceName, deviceType, s.ChannelId, sendDataOnlyOnChange, extension, functions, tags)
		if deviceTypeName, exist := config["deviceTypeName"]; exist {
			device.DeviceTypeName = deviceTypeName.(string)
		}
		devices[device.DeviceID] = device
	}
	s.ChannelCache.Devices = devices
	s.ChannelCache.Tags = tagMap
	s.ChannelCache.TagCount = tagCount
}

func (s *SimulateConnector) Init() {
	// 1. 初始化连接配置
	s.ConnectConfig = s.Config["connection"].(map[string]interface{})
	// 2. 解析配置为ChannelCache
	s.LoadConfig()

	reportConfig := &s.ReportConfig
	// 3. 解析配置得到采集频率任务
	reportConfig.Frequencies = make(map[int64][]connectors.Tag)
	channelFrequency := reportConfig.Frequency
	if reportConfig.CollectType == constants.Device {

		for deviceId, device := range s.ChannelCache.Devices {
			f := channelFrequency
			if device.Frequency > 0 {
				f = device.Frequency
			}
			if f > 0 {
				if tags, exist := reportConfig.Frequencies[f]; exist {
					tags = append(tags, device.Tags...)
				} else {
					reportConfig.Frequencies[f] = device.Tags
				}
			} else {

				cclog.Warn(fmt.Sprintf("device %v %v not set collect frequency", deviceId, device.DeviceName))
			}
		}
	} else if reportConfig.CollectType == constants.Tag {

		for _, device := range s.ChannelCache.Devices {
			f := channelFrequency
			if device.Frequency > 0 {
				f = device.Frequency
			}
			for _, tag := range device.Tags {
				if tag.Frequency > 0 {
					f = tag.Frequency
				}
				if f > 0 {
					if tags, exist := reportConfig.Frequencies[f]; exist {
						tags = append(tags, tag)
					} else {
						reportConfig.Frequencies[f] = make([]connectors.Tag, 8)
						reportConfig.Frequencies[f] = append(reportConfig.Frequencies[f], tag)
					}
				} else {

					cclog.Warn(fmt.Sprintf("device %v %v,Tag %v not set collect frequency", device.DeviceID, device.DeviceName, tag.TagId))
				}
			}
		}

	} else {

		if channelFrequency > 0 {

			for _, device := range s.ChannelCache.Devices {
				tags := s.ReportConfig.Frequencies[channelFrequency]
				if tags == nil {
					s.ReportConfig.Frequencies[channelFrequency] = device.Tags
				} else {
					s.ReportConfig.Frequencies[channelFrequency] = append(tags, device.Tags...)
				}

			}
		}
	}
}

func (s *SimulateConnector) Open() {

	// 3. 根据Connection 配置，建立南向连接
	s.Client = SimulateClient{ConnectionConfig: s.ConnectConfig}
	s.Client.Open(s.ChannelId)
	s.SetConnected(true, false)

	// 4. 注册设备
	s.ConnectorManager.RegisterDevices(lo.Values(s.ChannelCache.Devices))
	// 5. 建立采集任务

	for frequency, tags := range s.ReportConfig.Frequencies {

		_, err := s.Scheduler.Every(int(frequency)).Millisecond().Do(s.ReadTag, tags)
		if err != nil {
			cclog.Error(fmt.Sprintf("channal %v readTag failed.", s.ChannelName), err)
		}
	}
	s.Scheduler.StartAsync()
	// 6. 上线设备
	s.ConnectorManager.DeviceOnline(lo.Keys(s.ChannelCache.Devices))
}
func (s *SimulateConnector) ReadTag(tagSlice []connectors.Tag) {
	deviceTags := map[string][]connectors.Tag{}
	for _, tag := range tagSlice {
		if tags, exist := deviceTags[tag.DeviceId]; exist {
			deviceTags[tag.DeviceId] = append(tags, tag)
		} else {
			tags := make([]connectors.Tag, 0)
			deviceTags[tag.DeviceId] = append(tags, tag)
		}
	}
	for deviceId, tags := range deviceTags {
		values := s.Client.ReadTagValues(s.ChannelCache.Devices[deviceId], tags)
		event := events.SouthEvent{
			Type:       events.Telemetry,
			EventTopic: events.TelemetryTopic,
			DeviceId:   deviceId,
			Data:       values.Success,
		}
		s.ConnectorManager.SouthEventCh <- event
	}
}
func NewSimulatorConnector(connectorManager *connectors.ConnectorManager, channelId string, channelName string, config map[string]interface{}) connectors.Channel {
	var reportConfig connectors.ReportConfig
	if conf, ok := config[constants.ReportConfigKey]; ok {
		err := mapstructure.Decode(conf, &reportConfig)
		if err != nil {
			cclog.Error(fmt.Sprintf("decode report config failed,%v %v", channelName, conf), err)
		}
	}
	var channelInstance = connectors.ChannelInstance{ConnectorManager: connectorManager, ChannelId: channelId,
		ChannelName: channelName, Scheduler: gocron.NewScheduler(time.Local),
		ChannelCache: connectors.ChannelCache{}, ReportConfig: reportConfig, Connected: false, Enabled: true}
	simulator := &SimulateConnector{ChannelInstance: channelInstance, Config: config}
	simulator.ChannelCache = connectors.ChannelCache{
		Devices:   make(map[string]connectors.Device),
		Events:    nil,
		Extension: make(map[string]interface{}),
		TagCount:  0,
	}
	simulator.Init()
	cclog.Info(fmt.Sprintf("channel=%v success init", channelName))
	return simulator
}
