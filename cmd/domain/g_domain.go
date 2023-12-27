package domain

import (
	"github.com/go-co-op/gocron"
	mqtt "github.com/goiiot/libmqtt"
)

type MqttClientAddInfo struct {
	Address    string              `form:"address" json:"address" xml:"address" binding:"required"`
	ClientID   string              `form:"clientID" json:"clientID" xml:"clientID" binding:"required"`
	Username   string              `form:"username" json:"username" xml:"username" binding:"required"`
	Password   string              `form:"password" json:"password" xml:"password" binding:"required"`
	Keepalive  int64               `form:"keepalive" json:"keepalive" xml:"keepalive" binding:"required"`
	MockPolicy []PublishMockPolicy `form:"mockPolicy" json:"mockPolicy" xml:"mockPolicy"`
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
	clientID   string
	username   string
	password   string
	keepalive  int64
	scheduler  *gocron.Scheduler
	client     mqtt.Client
	enableMock bool
	MockPolicy []PublishMockPolicy `form:"mockPolicy" json:"mockPolicy" xml:"mockPolicy"`
}

func (G *GClientInfo) Scheduler() *gocron.Scheduler {
	return G.scheduler
}

func (G *GClientInfo) SetScheduler(scheduler *gocron.Scheduler) {
	G.scheduler = scheduler
}

func (G *GClientInfo) EnableMock() bool {
	return G.enableMock
}

func (G *GClientInfo) SetEnableMock(enableMock bool) {
	G.enableMock = enableMock
}

func (G *GClientInfo) ClientID() string {
	return G.clientID
}

func (G *GClientInfo) SetClientID(clientID string) {
	G.clientID = clientID
}

func (G *GClientInfo) Username() string {
	return G.username
}

func (G *GClientInfo) SetUsername(username string) {
	G.username = username
}

func (G *GClientInfo) Password() string {
	return G.password
}

func (G *GClientInfo) SetPassword(password string) {
	G.password = password
}

func (G *GClientInfo) Keepalive() int64 {
	return G.keepalive
}

func (G *GClientInfo) SetKeepalive(keepalive int64) {
	G.keepalive = keepalive
}

func (G *GClientInfo) Client() mqtt.Client {
	return G.client
}

func (G *GClientInfo) SetClient(client mqtt.Client) {
	G.client = client
}

func NewGClientInfo(server string, clientID string, username string, password string, keepalive int64) *GClientInfo {
	return &GClientInfo{Server: server, clientID: clientID, username: username, password: password, keepalive: keepalive, enableMock: true}
}

var ClientMaps = make(map[string]*GClientInfo, 8)
