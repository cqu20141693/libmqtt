package config

func GetStringOrDefault(key, defaultVal string) string {
	str := Viper.GetString(key)
	if str == "" {
		return defaultVal
	}
	return str
}
