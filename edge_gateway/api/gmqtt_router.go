package api

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

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

// CORSMiddleware 自定义跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许访问的域名，这里使用 * 表示允许所有域名访问
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		// 允许携带凭证（如 cookie）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func initGMqttApi() {

	// 跨域
	Router.Use(CORSMiddleware())

	// 静态文件
	frontPath := utils.GetStrEnv("FRONT_PATH", "/Users/gymd/work/ai-project/vue3-admin/dist")
	Router.Use(static.Serve("/", static.LocalFile(frontPath, false)))

	api := Router.Group("/api/gmqtt")

	// 设备api
	DeviceRoutes(api)
	ThingsModelRoutes(api)
	// 通道api
	ChannelRoutes(api)
	// 系统api
	SystemRoutes(api)
	// 用户api
	UserRoutes(api)
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
