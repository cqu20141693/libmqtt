package constants

type DataType = string

const (
	STRING  DataType = "String"
	INTEGER          = "Integer"
	LONG             = "Long"
	BOOL             = "Bool"
	Float            = "Float"
)

type TagType = string

const (
	TIMESERIES TagType = "telemetry"  //时序数据 只支持采集上报
	ATTRIBUTES         = "attributes" //属性数据 支持读取，写入和上报
)

type ReadWriteType uint8

const (
	ReadOnly  ReadWriteType = 1 //只支持读取
	WriteOnly               = 2 //只支持写入
	ReadWrite               = 3 //只支持写入
)

// TagQuality 点位质量
type TagQuality string

const (
	GOOD TagType = "GOOD" // 通信良好
	BAD          = "BAD"  // 通信断开
)
