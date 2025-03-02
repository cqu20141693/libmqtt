package enums

// BizType 业务枚举类型
type BizType int

// ThingsModelType 物模型类型
type ThingsModelType string

// 定义枚举值
const (
	Property ThingsModelType = "property" // iota 从 0 开始递增
	Command                  = "command"
)

// 定义枚举值
const (
	ChannelGroup BizType = iota // iota 从 0 开始递增
	Channel
	Device
)

// 为枚举类型实现 String() 方法
func (s BizType) String() string {
	switch s {
	case ChannelGroup:
		return "group"
	case Channel:
		return "channel"
	case Device:
		return "device"
	default:
		return "Unknown"
	}
}
