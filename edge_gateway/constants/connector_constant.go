package constants

const (
	ChannelConfigKey string = "server"
	ConnectConfigKey string = "connection"
	ReportConfigKey  string = "report"
)

const (
	FUNCTION       string = "function"
	FUNCTION_REPLY        = "function_reply"

	ATTRIBUTES_UPDATE       = "attributes_update"
	ATTRIBUTES_UPDATE_REPLY = "attributes_update_reply"

	ATTRIBUTES_READ       = "attributes_read"
	ATTRIBUTES_READ_REPLY = "attributes_read_reply"

	ADD_CHANNEL       = "add_channel"
	ADD_CHANNEL_REPLY = "add_channel_reply"

	STOP_CHANNEL       = "stop_channel"
	STOP_CHANNEL_REPLY = "stop_channel_reply"

	RESTART_CHANNEL       = "restart_channel"
	RESTART_CHANNEL_REPLY = "restart_channel_reply"
)

type CollectType = string

const (
	Channel CollectType = "channel"
	Device              = "device"
	Tag                 = "tag"
)
