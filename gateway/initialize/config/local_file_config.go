package config

import (
	"github.com/goiiot/libmqtt/gateway/initialize/logger/cclog"
	"github.com/spf13/viper"
)

var Viper = viper.New()

func init() {
	ReadLocalConfig("./resources/bootstrap.yml")

	active := Viper.GetString("will.profiles.active")
	if active != "" {
		ReadLocalConfig("./resources/bootstrap-" + active + ".yml")
	}
	//event.TriggerEvent(event.LocalConfigComplete)
}

func ReadLocalConfig(path string) {
	// 读取本地配置，支持相对路径
	Viper.SetConfigFile(path)
	Viper.SetConfigType("yaml")
	if err := Viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cclog.Debug("Config file not found; ignore error if desired")
		} else {
			cclog.Info("Config file was found but another error was produced. ", err)
		}
	}
}
