package cclog_test

import (
	"github.com/goiiot/libmqtt/cmd/initialize/logger/cclog"
	"go.uber.org/zap/zapcore"
	"testing"
)

func init() {
	cclog.SetLevel(zapcore.DebugLevel)
	//global.SetLogLevel(zapcore.DebugLevel)
}
func TestLog(t *testing.T) {
	//ccboot.Boot(nil)
	cclog.Debug("debug")
	cclog.Info("info")
	cclog.Warn("warn")
	cclog.Error("error")
}
