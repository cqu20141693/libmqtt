package crypto

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/domain"
	"github.com/tjfoc/gmsm/sm3"
	"strconv"
	"strings"
	"time"
)

var delimiter = ":"

func DoSignature(token string) (result string) {
	split := strings.Split(token, delimiter)
	if len(split) != 2 {
		return
	}
	switch split[0] {
	case domain.SM3:
		return sm3Sign(split, domain.SM3Login)
	case domain.SM4:
		return sm3Sign(split, domain.SM4Login)
	case domain.SHA256:
	case domain.AES:

	}

	return
}

func sm3Sign(split []string, loginType string) string {
	now := time.Now()
	timestamp := strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)
	nonce := utils.RandStr(8)
	join := strings.Join([]string{split[1], nonce, timestamp}, delimiter)
	hmacSM3 := doHmacSM3Base64(join, split[1])
	return strings.Join([]string{loginType, hmacSM3, nonce, timestamp}, delimiter)
}

func doSm3(data string) (sum []byte) {
	h := sm3.New()
	h.Write([]byte(data))
	sum = h.Sum(nil)
	return
}

func doSm3Base64(data string) (result string) {
	result = base64.StdEncoding.EncodeToString(doSm3(data))
	return
}

func doSm3HexStr(data string) (result string) {
	result = hex.EncodeToString(doSm3(data))
	return
}

func doHmacSM3(data, key string) (sum []byte) {
	mac := hmac.New(sm3.New, []byte(key))
	mac.Write([]byte(data))
	sum = mac.Sum(nil)
	return
}

func doHmacSM3Base64(data, key string) (result string) {
	result = base64.StdEncoding.EncodeToString(doHmacSM3(data, key))
	return
}

func doHmacSM3HexStr(data, key string) (result string) {
	result = hex.EncodeToString(doHmacSM3(data, key))
	return
}
