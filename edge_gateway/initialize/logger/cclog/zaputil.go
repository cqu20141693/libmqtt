package cclog

import (
	"github.com/goiiot/libmqtt/cmd/utils/file"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single log message.
type Entry struct {
	// Level is the severity of the log message.
	Level Level
	// Module is the dotted module name from the logger.
	Module string
	// Filename is the full path the file that logged the message.
	Filename string
	// Line is the line number of the Filename.
	Line int
	// Timestamp is when the log message was created
	Timestamp time.Time
	// Message is the formatted string from teh log call.
	Message string
	// Labels is the label associated with the log message.
	Labels []string
}

var loggoToZap = map[Level]zapcore.Level{
	TRACE:    zap.DebugLevel, // There's no zap equivalent to TRACE.
	DEBUG:    zap.DebugLevel,
	INFO:     zap.InfoLevel,
	WARNING:  zap.WarnLevel,
	ERROR:    zap.ErrorLevel,
	CRITICAL: zap.ErrorLevel, // There's no zap equivalent to CRITICAL.
}

// NewCCLogWriter returns a Writer that writes to the
// given zap logger.
func NewCCLogWriter(logger *zap.Logger) Writer {
	return zapCCLogWriter{
		logger: logger,
	}
}

// zapLoggoWriter implements a Writer by writing to a zap.Logger,
// so can be used as an adaptor from loggo to zap.
type zapCCLogWriter struct {
	logger *zap.Logger
}

// zapCCLogWriter implements Writer.Write by writing the entry
// to w.logger. It ignores entry.Timestamp because zap will affix its
// own timestamp.
func (w zapCCLogWriter) Write(entry Entry) {
	if ce := w.logger.Check(loggoToZap[entry.Level], entry.Message); ce != nil {
		ce.Write()
	}
}

func GetWriter(outputPath, name, rotateTime string, maxAge int64) (writer io.Writer, err error) {
	err = file.MkdirIfNecessary(outputPath)
	if err != nil {
		return
	}
	outputPath = outputPath + string(os.PathSeparator)
	rotateDuration, err := time.ParseDuration(rotateTime)
	writer, err = rotatelogs.New(filepath.Join(outputPath, name+"-%Y%m%d%H%M"),
		rotatelogs.WithRotationTime(rotateDuration), rotatelogs.WithMaxAge(time.Duration(maxAge)*rotateDuration),
		rotatelogs.WithLinkName(filepath.Join(outputPath, name)))
	return
}

var defaultLevel = zapcore.InfoLevel

// SetLevel 配置日志级别
func SetLevel(level zapcore.Level) {
	defaultLevel = level
}

func NewLogger(w io.Writer) *zap.Logger {
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(w),
		defaultLevel,
	)
	return zap.New(core)
}

func NewProductionLogger(w io.Writer) *zap.Logger {

	production, _ := zap.NewProduction()
	return production
}
