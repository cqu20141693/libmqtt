package utils

import "os"

// GetStrEnv 获取环境变量
func GetStrEnv(key string, defaultValue string) string {
	return GetEnv(key, defaultValue, func(s string) (string, error) {
		return s, nil
	})
}

// GetEnv 读取环境变量并返回指定类型的值
func GetEnv[T any](key string, defaultValue T, parseFunc func(string) (T, error)) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	if parseFunc == nil {
		return any(value).(T)
	}
	result, err := parseFunc(value)
	if err != nil {
		return defaultValue
	}

	return result
}
