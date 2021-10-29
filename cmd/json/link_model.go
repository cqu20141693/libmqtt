package json

type MqttInfo struct {
	ClientIdentifier string `json:clientIdentifier`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

type MqttAuthResult struct {
	Data    MqttInfo `json:data`
	Code    string   `json:code`
	Message string   `json:message`
}
