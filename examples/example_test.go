package examples

import (
	"fmt"
	"testing"
)

func TestClientConnect(t *testing.T) {
	ExampleClient()
}

// 定义枚举类型
type Status int

// 定义枚举值
const (
	StatusPending Status = iota // iota 从 0 开始递增
	StatusActive
	StatusCompleted
	StatusCancelled
)

// 为枚举类型实现 String() 方法
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "Pending"
	case StatusActive:
		return "Active"
	case StatusCompleted:
		return "Completed"
	case StatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func testEnum() {
	// 使用枚举值
	status := StatusActive

	// 打印枚举值及其字符串表示
	fmt.Printf("Status: %d, String: %s\n", status, status)

	// 枚举值比较
	if status == StatusActive {
		fmt.Println("The status is Active")
	}

	// 遍历枚举值
	fmt.Println("All status values:")
	for i := StatusPending; i <= StatusCancelled; i++ {
		fmt.Printf("%d: %s\n", i, i)
	}
}
func TestBasic(t *testing.T) {

	// 三元表达式 : 不支持，通过if,else

	// 枚举
	testEnum()

}
