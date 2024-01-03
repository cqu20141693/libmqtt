package utils

import (
	"fmt"
	"testing"
)

func TestTypeAssert(t *testing.T) {
	var m interface{}
	m = map[string]interface{}{
		"name": "tt",
		"age":  10,
	}
	m2 := m.(map[string]interface{})
	fmt.Printf("type assert success %v \n", m2)

	var x interface{}
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
