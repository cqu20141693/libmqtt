package examples

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {

	a := [3]int{1, 2, 3}
	b := [4]int{2, 4, 5, 6}
	fmt.Println(a, b)

	// 数组复制
	a1 := [...]int{1, 2, 3} // ... 会自动计算数组长度
	b1 := a1
	a1[0] = 100
	fmt.Println(a1, b1) // [100 2 3] [1 2 3]

	// 数组转切片
	ints := a[:2]
	fmt.Println(ints, len(ints), cap(ints))
}
func TestSliceInit(t *testing.T) {
	caps := 10
	size := 0
	// make 切片
	ids := make([]string, size, caps)
	ids2 := make([]string, size+1)

	printLenCap(ids)
	printLenCap(ids2)

	objs := []interface{}{1, 2, "3", 1.1, true}
	printInterfaceLenCap(objs)

	objs = append(objs, byte(1))
	printInterfaceLenCap(objs)
	var objs2 = make([]interface{}, len(objs))
	copy(objs2, objs)
	objs2 = append(objs2, "add")
	printInterfaceLenCap(objs2)
}

func TestSliceDelete(t *testing.T) {

}

func printLenCap[T comparable](nums []T) {
	fmt.Printf("len: %d, cap: %d %v\n", len(nums), cap(nums), nums)
}
func printInterfaceLenCap(nums []interface{}) {
	fmt.Printf("len: %d, cap: %d %v\n", len(nums), cap(nums), nums)
}
