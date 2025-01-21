package common

import (
	"github.com/goiiot/libmqtt/edge_gateway/random"
	"math/rand"
)

// MockInt
// 模拟整数
//
//	@param min
//	@param max
//	@return int
func MockInt(min int, max int) int {
	return min + (rand.Int() % (max - min))
}

func MockFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func MockBool() bool {
	return random.RandBool()
}
