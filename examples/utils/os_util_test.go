package utils

import (
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"strconv"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// 读取字符串类型的环境变量
	strValue := utils.GetEnv("MY_STRING_ENV", "default", func(s string) (string, error) {
		return s, nil
	})
	strValue = utils.GetStrEnv("MY_STRING_ENV", "default")
	fmt.Println("String value:", strValue)

	// 读取整数类型的环境变量
	intValue := utils.GetEnv("MY_INT_ENV", 42, strconv.Atoi)
	fmt.Println("Int value:", intValue)

	// 读取布尔类型的环境变量
	boolValue := utils.GetEnv("MY_BOOL_ENV", false, strconv.ParseBool)
	fmt.Println("Bool value:", boolValue)
}
