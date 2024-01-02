package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var SugarLogger *zap.SugaredLogger

func InitLogger() {

	productionConfig := zap.NewProductionConfig()
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	productionConfig.EncoderConfig = config
	logger, err := productionConfig.Build()
	if err != nil {
		return
	}
	SugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	productionConfig := zap.NewProductionConfig()
	config := zap.NewProductionEncoderConfig()
	productionConfig.EncoderConfig = config
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	return zapcore.NewJSONEncoder(config)
}

// getFileLogWriter  文件logger
func getFileLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}

var Router = gin.Default()

// GetRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func GetRoutes() {

	//api
	initGMqttApi()
}

func initGMqttApi() {
	api := Router.Group("/api/gmqtt")
	// 执行配置，建立client连接
	CreateClientRoutes(api)
	// 获取当前的客户端配置
	GetClients(api)
	//启用mock配置，开启定时推送任务
	StopMock(api)
	// 停止mock任务
	StartMock(api)
	// 更新mock任务
	UpdateMock(api)
	// 主动断开连接
	DisconnectClient(api)
	// 重连客户端
	ReconnectClient(api)
	// 一次性推送mqtt
	PublishMsg(api)

}
