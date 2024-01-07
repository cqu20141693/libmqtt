package domain

type ProtocolType uint8

const (
	// Simulator  模拟器
	Simulator ProtocolType = iota
	// Mqtt MQTT
	Mqtt
)

// 网关事件模型
// 通道事件
type AddChannelEvent struct {
	Protocol      string                 `json:"protocol"`
	ChannelId     string                 `json:"channelId"`
	ChannelName   string                 `json:"channelName"`
	ChannelConfig map[string]interface{} `json:"channelConfig"`
}
type UpdateChannelEvent struct {
	Protocol      string                 `json:"protocol"`
	ChannelId     string                 `json:"channelId"`
	ChannelName   string                 `json:"channelName"`
	ChannelConfig map[string]interface{} `json:"channelConfig"`
}

// 设备事件

// ota事件
