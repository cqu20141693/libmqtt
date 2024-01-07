package utils

func GetIntFromViper(config map[string]interface{}, key string) int {
	return int(config[key].(float64))
}
func GetInt64FromViper(config map[string]interface{}, key string) int64 {
	return int64(config[key].(float64))
}
