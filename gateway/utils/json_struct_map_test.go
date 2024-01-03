package utils

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
	"testing"
	"time"
)

type UserInfoVo struct {
	Id       string      `json:"id"`
	UserName string      `json:"user_name"`
	Address  []AddressVo `json:"address"`
}

type AddressVo struct {
	Address string `json:"address"`
}

var beforeMap = map[string]interface{}{
	"id":        "123",
	"user_name": "酒窝猪",
	"address":   []map[string]interface{}{{"address": "address01"}, {"address": "address02"}},
}

var User UserInfoVo

func init() {
	User = UserInfoVo{
		Id:       "01",
		UserName: "酒窝猪",
		Address: []AddressVo{
			{
				Address: "湖南",
			},
			{
				Address: "北京",
			},
		},
	}
}

// 性能高于TestMapToStructByJson
func TestMapToStructByMod(t *testing.T) {
	var afterStruct = UserInfoVo{}
	before := time.Now()
	err := mapstructure.Decode(beforeMap, &afterStruct)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result:%+v \n", time.Since(before))
	fmt.Printf("result:%+v \n", afterStruct)
}
func TestMapToStructByJson(t *testing.T) {
	beforeMap := map[string]interface{}{
		"id":        "123",
		"user_name": "酒窝猪",
		"address":   []map[string]interface{}{{"address": "address01"}, {"address": "address02"}},
	}
	var afterStruct = UserInfoVo{}
	before := time.Now()
	marshal, err := json.Marshal(beforeMap)
	if err != nil {
		fmt.Println("marshal:", err)
		return
	}
	err = json.Unmarshal(marshal, &afterStruct)
	if err != nil {
		fmt.Println("unmarshal:", err)
		return
	}
	fmt.Println(time.Since(before))
	fmt.Printf("resutlt: %+v", afterStruct)
}

func TestStructToMapByJson(t *testing.T) {
	var resultMap interface{}
	before := time.Now()
	jsonMarshal, _ := json.Marshal(User)
	err := json.Unmarshal(jsonMarshal, &resultMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(before))
	fmt.Printf("%+v", resultMap)
}

// 性能高于TestStructToMapByJson
func TestStructToMapByReflect(t *testing.T) {
	var resultMap = make(map[string]interface{}, 10)
	before := time.Now()

	ty := reflect.TypeOf(User)
	v := reflect.ValueOf(User)
	for i := 0; i < v.NumField(); i++ {
		resultMap[strings.ToLower(ty.Field(i).Name)] = v.Field(i).Interface()
	}
	fmt.Println(time.Since(before))
	fmt.Printf("%+v", resultMap)
}
