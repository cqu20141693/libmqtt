package domain

import mqtt "github.com/goiiot/libmqtt"

const (
	LineStart = "> "
)

// signature algorithm
const (
	SM3    = "SM3"
	SHA256 = "SHA256"
)

// crypt algorithm
const (
	SM4 = "SM4"
	AES = "AES"
)

// login type
const (
	SM3Login    = "G-HmacSM3"
	SM4Login    = "GC-HmacSM3"
	ShA256Login = "G-HmacSHA256"
	AESLogin    = "GC-HmacSHA256"
)

// data encoding type
const (
	Int    = "int"
	Long   = "long"
	Float  = "float"
	Double = "double"
	String = "string"
	Json   = "json"
	Bin    = "bin"
)

var IdGenerator int64

// ClientMap clientId ->*mqtt.AsyncClient
var ClientMap = make(map[string]*mqtt.AsyncClient, 8)
var CmdMap = make(map[string]*ClientInfo, 8)
var ClientInfoMap = make(map[mqtt.Client]*ClientInfo, 8)

type ClientInfo struct {
	Cmd      string
	Server   string
	ClientID string
	Username string
	Model    string
	Token    string
	GK       string
	Extend   []string
	Id       string
	welcome  WelcomeInfo
}

func (c *ClientInfo) Welcome() WelcomeInfo {
	return c.welcome
}

func (c *ClientInfo) SetWelcome(welcome WelcomeInfo) {
	c.welcome = welcome
}

func (c *ClientInfo) SetId(id string) {
	c.Id = id
}

func NewClientInfo(cmd, server, clientID, username, model, token, gk string, extend []string) *ClientInfo {
	return &ClientInfo{Cmd: cmd, Server: server, ClientID: clientID, Username: username, Model: model, Token: token, GK: gk, Extend: extend}
}

type WelcomeInfo struct {
	Info         string `json:"info"`
	GK           string `json:"gk"`
	LinkAddress  string `json:"linkAddress"`
	CryptoSecret string `json:"cryptoSecret"`
}
