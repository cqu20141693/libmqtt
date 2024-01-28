package connectors

import (
	"github.com/go-co-op/gocron"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/events"
	"github.com/samber/lo"
	"time"
)

//	Converter
//
// 南北向数据转换器
type Converter interface {
	UpConvert(data interface{}) events.SouthEventType
	DownConvert(event events.NorthCmdEvent) interface{}
}

// Channel 定义通道功能
type Channel interface {
	// Init
	//初始化连接配置
	// 初始化ChannelCache
	// 如果存在 Converter ，初始化Converter
	Init()
	//Open
	//建立南向设备连接
	//初始化采集任务
	//解析配置得到采集频率任务
	Open()

	// Close 关闭通道（主动关闭，和连接失败关闭）
	//  @param remote: 是否为用户主动关闭 区别在于用户主动关闭时，不需要打印错误日志
	// 需要关闭南向设备连接
	// 停止数据采集任务
	Close(remote bool)

	// IsConnected 获取通道是否已连接
	//  @return bool 连接状态
	IsConnected() bool

	// GetId 获取通道ID
	//  @return string
	GetId() string

	//ServerSideRpcHandler 北向rpc调用 异步调用
	ServerSideRpcHandler(eventData RpcEvent)

	//ServerSideRpcHandlerSync 北向rpc调用 同步调用
	//  @return map[string]ReadTagResult deviceId->result
	ServerSideRpcHandlerSync(eventData RpcEvent) RpcResult

	// AttributesUpdate 北向调用属性更新,异步读
	//  @param eventData
	AttributesUpdate(eventData UpdateTagEvent)

	// AttributesUpdateSync 北向调用属性更新,异步读
	//  @param eventData
	//  @return map[string]ReadTagResult deviceId->result
	AttributesUpdateSync(eventData UpdateTagEvent) UpdateTagResult

	// AttributesRead  北向调用属性读取,异步读
	//  @param eventData
	AttributesRead(eventData ReadTagEvent)

	// AttributesReadSync 北向调用属性读取,异步读
	//  @param eventData
	//  @return map[string]ReadTagResult deviceId->result
	AttributesReadSync(eventData ReadTagEvent) ReadTagResult
}

// ChannelInstance 通道父类
type ChannelInstance struct {
	ConnectorManager *ConnectorManager
	ChannelId        string
	ChannelName      string
	Converter        Converter         // 上行数据转换器
	Scheduler        *gocron.Scheduler //任务调度器
	ChannelCache     ChannelCache
	ReportConfig     ReportConfig // 采集配置
	Connected        bool         // 是否连接成功
	Enabled          bool         // 是否启用
}

func (c *ChannelInstance) SetConnected(connected bool, deviceOnline bool) {

	if connected == c.Connected {
		// 通道状态未变化
		return
	} else {
		c.Connected = connected
		if c.Connected {
			if deviceOnline {
				c.ConnectorManager.DeviceOnline(lo.Keys(c.ChannelCache.Devices))
			}
		} else {
			c.ConnectorManager.DeviceOffline(lo.Keys(c.ChannelCache.Devices))
		}
	}
}

func (c *ChannelInstance) PushSouthEvent(event events.SouthEvent) {
	c.ConnectorManager.SouthEventCh <- event
}

type ReportConfig struct {
	// SendDataOnlyOnChange true 变化上报，否则采集上报
	// 目前未实现，每个采集点位都已配置
	// 可以考虑前端创建设备时使用改值作为默认值
	SendDataOnlyOnChange bool                  `json:"sendDataOnlyOnChange"`
	CollectType          constants.CollectType `json:"collectType"`
	// 采集频率
	Frequency int64 `json:"frequency"`
	// 多频下采集频率
	// 多频下存在同一个设备一个频率，或者同一设备下存在不同频率
	Frequencies map[int64][]Tag
}

type Function struct {
	Method    string                 `json:"method"`
	Name      string                 `json:"name"`
	DeviceID  string                 `json:"device_id"`
	Extension map[string]interface{} `json:"extension"`
}

// Device 设备模型
type Device struct {
	DeviceID       string `json:"deviceId"`
	DeviceName     string `json:"deviceName"`
	DeviceType     string `json:"deviceType"`
	DeviceTypeName string `json:"deviceTypeName"`
	ChannelID      string `json:"channelID"`
	// SendDataOnlyOnChange true 变化上报，否则采集上报
	// 可以考虑前端创建点位时默认值使用设备配置
	SendDataOnlyOnChange bool  `json:"sendDataOnlyOnChange"`
	Frequency            int64 `json:"frequency"`
	// Extension  设备扩展，不同的设备扩展可能不一样
	Extension map[string]interface{} `json:"extension"`
	// Attributes 可读属性
	Attributes []string `json:"attributes"`
	// AttributeUpdates 可写属性
	AttributeUpdates []string `json:"attributeUpdates"`
	// 时许属性
	Timeseries  []string   `json:"timeseries"`
	Functions   []Function `json:"functions"`
	Tags        []Tag      `json:"tags"`
	tagMap      map[string]Tag
	functionMap map[string]Function
}

// NewDevice 创建设备模型
//
//	@param deviceID
//	@param deviceName
//	@param deviceType
//	@param channelID
//	@param sendDataOnlyOnChange
//	@param config
//	@param attributes
//	@param attributesUpdates
//	@param timeseries
//	@param function
//	@return *Device
func NewDevice(deviceID string, deviceName string, deviceType string, channelID string, sendDataOnlyOnChange bool, extension map[string]interface{}, functions []Function, tags []Tag) Device {
	device := Device{
		DeviceID: deviceID, DeviceName: deviceName,
		DeviceType: deviceType, ChannelID: channelID,
		SendDataOnlyOnChange: sendDataOnlyOnChange, Extension: extension,
		Functions: functions, Tags: tags}
	device.init()
	return device
}

// init 构建设备tagMap
//
//	@receiver d
func (d *Device) init() {
	d.tagMap = map[string]Tag{}
	d.functionMap = map[string]Function{}
	if d.Tags != nil {
		for _, tag := range d.Tags {
			d.tagMap[tag.TagId] = tag
		}
	}
	if d.Functions != nil {
		for _, function := range d.Functions {
			d.functionMap[function.Method] = function
		}
	}

}

// FindTag 查询设备Tag
//
//	@receiver d
//	@param tagId
//	@return Tag
func (d *Device) FindTag(tagId string) Tag {
	return d.tagMap[tagId]
}

// FindFunction 查询功能
//
//	@receiver d
//	@param fId
//	@return Function
func (d Device) FindFunction(fId string) Function {
	return d.functionMap[fId]
}

// Tag  点位模型
type Tag struct {
	Name     string             `json:"name"`
	TagId    string             `json:"tagId"`
	DeviceId string             `json:"deviceId"`
	Value    interface{}        `json:"value"`
	DataType constants.DataType `json:"DataType"`
	//  ReadWrite TagType= ATTRIBUTES 时生效
	ReadWrite constants.ReadWriteType `json:"readWrite"`
	TagType   constants.TagType       `json:"tagType"`
	Quality   constants.TagQuality    `json:"quality"`
	Frequency int64                   `json:"frequency"`
	// SendDataOnlyOnChange true 变化上报，否则采集上报
	SendDataOnlyOnChange bool                   `json:"sendDataOnlyOnChange"`
	Extension            map[string]interface{} `json:"extension"`
}

func NewTag(name string, tagId string, deviceId string, dataType constants.DataType, readWrite constants.ReadWriteType, tagType constants.TagType, frequency int64, sendDataOnlyOnChange bool, extension map[string]interface{}) Tag {
	return Tag{Name: name, TagId: tagId, DataType: dataType, DeviceId: deviceId, ReadWrite: readWrite, TagType: tagType, Frequency: frequency,
		SendDataOnlyOnChange: sendDataOnlyOnChange, Extension: extension, Quality: constants.BAD}
}

type EventCode string

// Event 事件模型
type Event struct {
	ChannelId   string
	ChannelName string
	Code        EventCode
	Message     string
	Time        time.Time
}

// ChannelCache 通道缓存信息
type ChannelCache struct {
	Devices   map[string]Device      //设备信息
	Tags      map[string]Tag         //设备信息
	Events    map[string]Event       //事件信息
	Extension map[string]interface{} // 扩展信息
	TagCount  uint16                 //点位数量
}
type ReadTagEvent struct {
	DeviceId string
	TagIds   []string
}
type RpcEvent struct {
	DeviceId string
	Method   string
	Params   interface{}
}
type UpdateTagEvent struct {
	DeviceId string
	TagValue map[string]interface{}
}
type ReadTagResult struct {
	// 点位的时间戳
	Ts       int64
	DeviceId string
	// Fails 读取失败点位
	// tagId -> reason
	Fails map[string]string
	//  Success 读取成功读点诶
	// tagId->value
	Success map[string]interface{}
}
type UpdateTagResult struct {
	DeviceId string
	// Fails 读取失败点位
	// tagId -> reason
	Fails map[string]string
	//  Success 读取成功读点诶
	// tagId->value
	Success []string
}
type RpcResult struct {
	DeviceId string
	//  Success  true 成功，否则失败
	Success bool
	// Fail 失败原因
	FailReason string
	// 请求数据
	Data interface{}
}
