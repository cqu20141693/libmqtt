package domain

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	mqtt "github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
)

type MqttClientAddInfo struct {
	Address         string              `form:"address" json:"address" xml:"address" binding:"required"`
	ClientID        string              `form:"clientID" json:"clientID" xml:"clientID" binding:"required"`
	Username        string              `form:"username" json:"username" xml:"username" binding:"required"`
	Password        string              `form:"password" json:"password" xml:"password" binding:"required"`
	Keepalive       int64               `form:"keepalive" json:"keepalive" xml:"keepalive" binding:"required"`
	MockPolicy      []PublishMockPolicy `form:"mockPolicy" json:"mockPolicy" xml:"mockPolicy"`
	ProtocolVersion byte
	client          mqtt.Client
}

func (m *MqttClientAddInfo) PrintClientMetric() {
	client := m.client
	if client != nil && client.PubMetric != nil {
		cclog.SugarLogger.Info(fmt.Sprintf("device: %s, publish: %d", m.ClientID, client.PubMetric.Count()))
	}
}

func (m *MqttClientAddInfo) SetClient(client mqtt.Client) {
	m.client = client
}

func NewMqttClientAddInfoWithVersion(address string, clientID string, username string, password string, keepalive int64, version byte) *MqttClientAddInfo {
	return &MqttClientAddInfo{Address: address, ClientID: clientID, Username: username, Password: password, Keepalive: keepalive, ProtocolVersion: version}
}

func NewMqttClientAddInfo(address string, clientID string, username string, password string, keepalive int64) *MqttClientAddInfo {
	return NewMqttClientAddInfoWithVersion(address, clientID, username, password, keepalive, 4)
}

type PublishMockPolicy struct {
	Enable bool   `form:"enable" json:"enable" xml:"enable" binding:"required"`
	Topic  string `form:"topic" json:"topic" xml:"topic" binding:"required"`
	Qos    byte   `form:"qos" json:"qos" xml:"qos" binding:"required"`
	// 频率 毫秒
	Frequency int    `form:"frequency" json:"frequency" xml:"frequency" binding:"required"`
	Message   string `form:"message" json:"message" xml:"message" binding:"required"`
}

// 广 域mqtt domain

type GClientInfo struct {
	Server     string
	ClientID   string
	Username   string
	Password   string
	Keepalive  int64
	Scheduler  gocron.Scheduler
	Client     mqtt.Client
	EnableMock bool
	Connected  bool
	MockPolicy []PublishMockPolicy `form:"mockPolicy" json:"mockPolicy" xml:"mockPolicy"`
}

func (G *GClientInfo) GetScheduler() gocron.Scheduler {
	return G.Scheduler
}

func (G *GClientInfo) SetScheduler(scheduler gocron.Scheduler) {
	G.Scheduler = scheduler
}

func (G *GClientInfo) GetEnableMock() bool {
	return G.EnableMock
}

func (G *GClientInfo) SetEnableMock(enableMock bool) {
	G.EnableMock = enableMock
}

func (G *GClientInfo) GetClientID() string {
	return G.ClientID
}

func (G *GClientInfo) SetClientID(clientID string) {
	G.ClientID = clientID
}

func (G *GClientInfo) GetUsername() string {
	return G.Username
}

func (G *GClientInfo) SetUsername(username string) {
	G.Username = username
}

func (G *GClientInfo) GetPassword() string {
	return G.Password
}

func (G *GClientInfo) SetPassword(password string) {
	G.Password = password
}

func (G *GClientInfo) GetKeepalive() int64 {
	return G.Keepalive
}

func (G *GClientInfo) SetKeepalive(keepalive int64) {
	G.Keepalive = keepalive
}

func (G *GClientInfo) GetClient() mqtt.Client {
	return G.Client
}

func (G *GClientInfo) SetClient(client mqtt.Client) {
	G.Client = client
}

func NewGClientInfo(server string, clientID string, username string, password string, keepalive int64) *GClientInfo {
	return &GClientInfo{Server: server, ClientID: clientID, Username: username, Password: password, Keepalive: keepalive, EnableMock: true}
}

var ClientMaps = make(map[string]*GClientInfo, 8)
