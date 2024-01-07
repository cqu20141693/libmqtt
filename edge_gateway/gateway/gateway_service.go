package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/connectors"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/edge_gateway/platform"
	"github.com/goiiot/libmqtt/edge_gateway/register"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var Service *GatewayService
var Manger *connectors.ConnectorManager

func StartGateway() {

	northCh := make(chan events.NorthEvent, 1)
	southCh := make(chan events.SouthEvent, 1)
	configPath := viper.GetString("gateway.configPath")
	if configPath == "" {
		// 获取当前文件的路径
		_, filename, _, _ := runtime.Caller(0)
		root := path.Dir(path.Dir(path.Dir(filename)))
		configPath = root + "/resources/channels"
	}
	Service = NewGatewayService(northCh, southCh, configPath)
	Manger = &connectors.ConnectorManager{
		NorthEventCh:     northCh,
		SouthEventCh:     southCh,
		ChannelInstances: map[string]connectors.Channel{},
		Devices:          map[string]string{},
	}
	// 启动Manager
	Manger.Start()
	// 启动网关
	Service.start()
}

// GatewayService 网关服务
type GatewayService struct {

	// NorthCh 北向平台事件
	// 事件需要保持顺序行
	NorthCh chan events.NorthEvent

	//  SouthCh
	// 事件需要保持顺序行
	SouthCh chan events.SouthEvent

	//  ChannelConfigPath	 配置文件夹
	ChannelConfigPath string
}

func NewDefaultGatewayService() *GatewayService {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(filename)))
	dir := root + "/resources/channels"
	return &GatewayService{NorthCh: make(chan events.NorthEvent, 1), SouthCh: make(chan events.SouthEvent), ChannelConfigPath: dir}
}
func NewGatewayService(northCh chan events.NorthEvent, southCh chan events.SouthEvent, channelConfigPath string) *GatewayService {
	return &GatewayService{NorthCh: northCh, SouthCh: southCh, ChannelConfigPath: channelConfigPath}
}

func (receiver *GatewayService) start() {

	// 分发南向事件
	go DispatchSouthEvent(receiver)
	// 读取通道配置
	// 调用 filepath.Walk 函数遍历文件夹
	var configs []map[string]interface{}
	configs = receiver.ReadConfigDir(configs)
	//cclog.Info("load channel config,prepare send AddChannelEvent.", configs)
	// 启动通道
	for _, config := range configs {
		serverConfig := config["server"]
		type Server struct {
			Name      string                 `json:"name"`
			Type      string                 `json:"type"`
			Dynamic   bool                   `json:"dynamic"`
			ChannelId string                 `json:"channelId"`
			Extension map[string]interface{} `json:"extension"`
		}
		var server Server
		err := mapstructure.Decode(serverConfig, &server)
		if err != nil {
			sprintf := fmt.Sprintf("decode serverConfig faield %v", config)
			cclog.Error(sprintf, err)
		}
		if generator, ok := register.DefaultConnectors[server.Type]; ok {

			channel := generator(Manger, server.ChannelId, server.Name, serverConfig.(map[string]interface{}))
			channel.Open()
			Manger.ChannelInstances[channel.GetId()] = channel
		} else {

			cclog.Warn(fmt.Sprintf("not support connector %v", server))
		}
	}
}

func DispatchSouthEvent(receiver *GatewayService) {
	for true {

		select {
		case event, more := <-receiver.SouthCh:
			if !more {
				return
			}
			for _, platformClient := range platform.PlatformClientMaps {
				switch event.Type {
				case events.Event:
					platformClient.PublishEvent(platform.Event{
						EventTopic: event.EventTopic,
						DeviceId:   event.DeviceId,
						Data:       event.Data,
					})
				case events.Telemetry:
					platformClient.PublishTelemetry(platform.Telemetry{
						EventTopic: event.EventTopic,
						DeviceId:   event.DeviceId,
						Data:       event.Data,
					})
				case events.DeviceEvent:
					platformClient.PublishDeviceEvent(platform.DeviceEvent{
						EventTopic: event.EventTopic,
						DeviceId:   event.DeviceId,
						Data:       event.Data,
					})
				case events.Reply:
					platformClient.PublishReply(platform.Reply{
						EventTopic: event.EventTopic,
						DeviceId:   event.DeviceId,
						Data:       event.Data,
					})
				}
			}
		}
	}
}

func (receiver *GatewayService) ReadConfigDir(configs []map[string]interface{}) []map[string]interface{} {
	err := filepath.Walk(receiver.ChannelConfigPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断文件是否为 JSON 文件
		if filepath.Ext(path) == ".json" {
			// 读取 JSON 文件内容
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// 解析 JSON 数据
			var jsonData map[string]interface{}
			err = json.Unmarshal(data, &jsonData)
			if err != nil {
				return err
			}
			configs = append(configs, jsonData)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return configs
}
