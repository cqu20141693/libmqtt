package events

type TelemetryData struct {
	DeviceName  string
	DeviceType  string
	DeviceId    string
	Telemetries []map[string]interface{}
	Attributes  []map[string]interface{}
	Uuid        string
}

type TelemetryEvent struct {
	ChannelId string
	Data      TelemetryData
}

// NorthCmdEvent 北向功能调用
type NorthCmdEvent struct {
	ChannelId string
	Command   interface{}
	CmdType   string
}

type NorthChannelEvent struct {
	ChannelId   string
	ChannelName string
	Protocol    string
	Config      map[string]interface{}
}

// NorthEvent 北向平台事件
type NorthEvent struct {
	Type string
	Data interface{}
}

type SouthEventType string

const (
	Telemetry   SouthEventType = "telemetry"
	DeviceEvent SouthEventType = "deviceEvent"
	Event       SouthEventType = "event"
	Reply       SouthEventType = "reply"
)

type EventTopicType string

const (
	OnlineTopic    EventTopicType = "online"
	OfflineTopic   EventTopicType = "offline"
	RegisterTopic  EventTopicType = "register"
	TelemetryTopic EventTopicType = "telemetry"
	HealthTopic    EventTopicType = "health"
)

// SouthEvent 南向设备事件
type SouthEvent struct {
	Type       SouthEventType
	EventTopic EventTopicType
	DeviceId   string
	Data       interface{}
}
