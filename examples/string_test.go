package examples

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	var s string = "{\"aType\":\"sysLog\",\"className\":\"Metrics\",\"createTime\":1704761970000,\"id\":\"1744388577479671808\",\"level\":\"INFO\",\"lineNumber\":44,\"message\":\"{}\",\"methodName\":\"init\",\"name\":\"Metrics\",\"threadId\":\"546\",\"threadName\":\"1\"}"
	fmt.Println(s)
}
