package examples

import (
	"fmt"
	"strconv"
	"testing"
)

func TestTypeAssert(t *testing.T) {
	var m interface{}
	m = map[string]interface{}{
		"name": "tt",
		"age":  10,
	}
	// 类型断言
	m2, ok := m.(map[string]interface{})
	if ok {
		fmt.Printf("type assert success %v \n", m2)
	}

	var x interface{}
	// 类型
	typeSwitch(x)
	x = 1
	typeSwitch(x)
	x = 1.01
	typeSwitch(x)
	x = true
	typeSwitch(x)
	x = "你好"
	typeSwitch(x)
	x = func(i int) float64 {
		return float64(i)
	}
	typeSwitch(x)
	x = map[string]interface{}{
		"name": "tt",
		"age":  10,
	}
	typeSwitch(x)
}

func typeSwitch(x interface{}) {
	switch i := x.(type) {
	case nil:
		fmt.Printf("x is nil\n") // type of i is type of x (interface{})
	case int:
		fmt.Printf("x is int %v \n", i) // type of i is int
	case float64:
		fmt.Printf("x is float64 %v\n", i) // type of i is float64
	case func(int) float64:
		fmt.Printf("x is func(int) float64 %v\n", x) // type of i is func(int) float64
	case bool, string:
		fmt.Printf("type is bool or string %v \n", i) // type of i is type of x (interface{})
	default:
		fmt.Printf("don't know the type %v \n", i) // type of i is type of x (interface{})
	}
}

func TestTypeConvert(t *testing.T) {
	// 测试字符串转其他类型
	testStrconv()
}

func testStrconv() {
	var i int = 1
	str := strconv.Itoa(i)
	_, err := strconv.Atoi(str)
	if err != nil {
		return
	}
	fmt.Println("strconv Atoi Itoa success")
}
