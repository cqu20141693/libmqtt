package domain

// Gateway 网关配置模型
type Gateway struct {
	AutoConnect bool        `yaml:"autoConnect" json:"autoConnect"`
	Version     string      `yaml:"version" json:"version"`
	HealthCheck HealthCheck `yaml:"healthCheck" json:"healthCheck"`
	Platform    Platform    `yaml:"platform" json:"platform"`
	IsUpdate    bool        `yaml:"isUpdate" json:"isUpdate"`
	RestPort    int         `yaml:"restPort" json:"restPort"`
}

type HealthCheck struct {
	Enabled        bool `yaml:"enabled" json:"enabled"`
	ReportInterval int  `yaml:"reportInterval" json:"reportInterval"`
}

type Platform struct {
	Qos       int    `yaml:"qos" json:"qos"`
	Port      int    `yaml:"port" json:"port"`
	Username  string `yaml:"username" json:"username"`
	ClientId  string `yaml:"clientId" json:"clientId"`
	Password  string `yaml:"password" json:"password"`
	Host      string `yaml:"host" json:"host"`
	KeepAlive int64  `yaml:"keepAlive" json:"keepAlive"`
	MaxSize   int    `yaml:"maxSize" json:"maxSize"`
}
