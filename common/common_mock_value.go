package common

import "math/rand"

// MockInt
// 模拟整数
//
//	@param min
//	@param max
//	@return int
func MockInt(min int, max int) int {
	return min + (rand.Int() % (max - min))
}
