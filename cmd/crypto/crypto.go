package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
)

func DoSM4(data, key, iv []byte) (ciphertext []byte, err error) {
	return sm4Encrypt(data, key, iv)
}
func DoSM4Decrypt(data, key, iv []byte) (plaintText []byte, err error) {

	return sm4Decrypt(data, key, iv)
}

// 明文数据填充
// pkcs5填充
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// 去掉明文后面的填充数据
func unpaddingLastGroup(plainText []byte) []byte {
	//1.拿到切片中的最后一个字节
	length := len(plainText)
	lastChar := plainText[length-1]
	//2.将最后一个数据转换为整数
	number := int(lastChar)
	return plainText[:length-number]
}

func sm4Encrypt(plainText, key, iv []byte) (cipherText []byte, err error) {
	if len(iv) < 16 {
		return nil, errors.New("iv length less  16")
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData := pkcs5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText = make([]byte, len(origData))
	blockMode.CryptBlocks(cipherText, origData)
	return
}

func sm4Decrypt(cipherText, key, iv []byte) (plainText []byte, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			fmt.Println(err1)
			err = errors.New("sm4 decrypt failed")
		}
	}()
	if len(iv) < 16 {
		return nil, errors.New("iv length less  16")
	}
	block, err1 := sm4.NewCipher(key)
	if err1 != nil {
		return nil, err1
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:16])
	plainText = make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = unpaddingLastGroup(plainText)
	return
}

func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
