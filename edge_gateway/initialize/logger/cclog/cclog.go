package cclog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

var logs = map[string]Writer{"console": getConsoleWriter()}
var console Writer

var SugarLogger *zap.SugaredLogger

func init() {

	logger, done := getLogger()
	if done {
		return
	}
	SugarLogger = logger.Sugar()
}

func getLogger() (*zap.Logger, bool) {
	productionConfig := zap.NewProductionConfig()
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	productionConfig.EncoderConfig = config
	logger, err := productionConfig.Build()
	if err != nil {
		return nil, true
	}
	return logger, false
}

func getConsoleWriter() Writer {
	logger, done := getLogger()
	if done {
		return nil
	}
	console = NewCCLogWriter(logger)
	return console
}
func newLogger(w io.Writer, level zapcore.Level) *zap.Logger {
	eConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(eConfig),
		zapcore.AddSync(w),
		level,
	)
	return zap.New(core, zap.AddCaller())
}

func AddWriter(key string, writer Writer) {
	logs[key] = writer
}
func AddLogger(key string, level zapcore.Level, writer io.Writer) {
	logs[key] = NewCCLogWriter(newLogger(writer, level))
}

func DebugEntry(msg string) Entry {
	return Entry{Level: DEBUG, Message: msg, Timestamp: time.Now()}
}
func InfoEntry(msg string) Entry {
	return Entry{Level: INFO, Message: msg, Timestamp: time.Now()}
}
func WarnEntry(msg string) Entry {
	return Entry{Level: WARNING, Message: msg, Timestamp: time.Now()}
}
func ErrorEntry(msg string) Entry {
	return Entry{Level: ERROR, Message: msg, Timestamp: time.Now()}
}

func Info(v ...any) {
	for _, log := range logs {
		log.Write(InfoEntry(fmt.Sprint(v...)))
	}
}
func Debug(v ...any) {
	for _, log := range logs {
		log.Write(DebugEntry(fmt.Sprint(v...)))
	}
}
func Warn(v ...any) {
	for _, log := range logs {
		log.Write(WarnEntry(fmt.Sprint(v...)))
	}
}
func Error(v ...any) {
	for _, log := range logs {
		log.Write(ErrorEntry(fmt.Sprint(v...)))
	}
}
