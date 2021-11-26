package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
	"hash"
	"testing"
)

func TestSM4(t *testing.T) {
	src := []byte("go-data")
	data := []byte("java-data")
	key := []byte("2fw6oC2eVtDKraks")
	cipherText, err := DoSM4(src, key, key)
	if err != nil {
		panic(err)
	}
	plainText, err := DoSM4Decrypt(cipherText, key, key)
	if err != nil {
		panic(err)
	}
	flag := bytes.Equal(src, plainText)
	fmt.Printf("SM4是否解密成功：%v:%s \n", flag, base64.StdEncoding.EncodeToString(cipherText))

	_, err = base64.StdEncoding.Decode(cipherText, []byte("vfXMqpJaDxRM4AGNI615LA=="))
	if err != nil {
		return
	}
	doSM4, _ := DoSM4Decrypt(cipherText, key, key)
	flag = bytes.Equal(data, doSM4)
	fmt.Printf("SM4 解密成功:%v:%s \n", flag, string(doSM4))
}

func TestDoSignature(t *testing.T) {
	signature := DoSignature("GSM3:2fw6oC2eVtDKraks")
	fmt.Printf("sm3 signature value:%s", signature)
}

func TestSM3(t *testing.T) {
	data := "test"
	key := "key"
	h := sm3.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	fmt.Printf("digest value is: %x\n", sum)
	messageMAC := []byte("02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159")
	if ok, bytes := ValidMAC(sm3.New, []byte(data), messageMAC, []byte(key)); !ok {
		toString := hex.EncodeToString(bytes)
		fmt.Printf("hmacSha256 digest valid failed,expect value is:%x %s\n", bytes, toString)
	} else {
		fmt.Printf("hmacSha256 digest valid success")
	}
	if ok, bytes := ValidMAC(sha256.New, []byte(data), messageMAC, []byte(key)); !ok {
		fmt.Printf("hmacSha256 digest valid failed,expect value is:%x\n", bytes)
	} else {
		fmt.Printf("hmacSha256 digest valid success")
	}

}

func ValidMAC(h func() hash.Hash, message, messageMAC, key []byte) (bool, []byte) {
	mac := hmac.New(h, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC), expectedMAC
}

func TestBase64(t *testing.T) {
	data := "zW9/JDYsfrqQ72lHUaCuXRMOhIjZZSPbTv+9BToJ1qJdMYaQ7yxBSQrmsZuEym50TjMzwOR0bIIGbK/vtHsIgmRC3DtXW1hjqSobIvQMCRHxnA/jXJDfxpR1rnaHDwQ0gXfCrvlqrDyvaM8nMx5MTpX5sKWnXJTWuvidGPSbsc8IISMkdp8FIJkAviFwZOWvlklyScYeByOvStgIijpkfg=="
	text, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}
	fmt.Println(string(text))
}
