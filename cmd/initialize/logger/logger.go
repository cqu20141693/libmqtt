package logger

import (
	"github.com/goiiot/libmqtt/cmd/initialize/config"
	"github.com/goiiot/libmqtt/cmd/initialize/logger/cclog"
	"github.com/goiiot/libmqtt/cmd/utils/file"
	"os"
	"strings"
)

func init() {
	ConfigRotate()
}
func ConfigRotate() {
	cclog.AddWriter("console", cclog.NewCCLogWriter(cclog.NewLogger(os.Stdout)))
	addDefaultConfig()
	rotateTime := config.Viper.GetString("will.log.rotate-time")
	maxAge := config.Viper.GetInt64("will.log.max-age")
	var path string
	if logDir := config.GetStringOrDefault("will.log.dir", ""); logDir != "" {
		if strings.Contains(logDir, "/") {
			path = logDir
		} else {
			path = file.GetCurrentPath() + string(os.PathSeparator) + "log"
		}
	} else {
		path = file.GetCurrentPath()
	}

	service := config.GetStringOrDefault("will.application.name", "witi-service")
	writer, err := cclog.GetWriter(path, service+".log", rotateTime, maxAge)
	if err != nil {
		cclog.Error("rotate writer create failed")
		return
	}
	rotate := cclog.NewCCLogWriter(cclog.NewLogger(writer))
	cclog.AddWriter("rotate", rotate)
	cclog.Debug("trigger log complete")
}
func addDefaultConfig() {
	config.Viper.Set("will.log.max-age", 3)
	config.Viper.Set("will.log.rotate-time", "24h")
}
