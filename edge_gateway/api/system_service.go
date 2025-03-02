package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// SystemInfo 定义系统信息结构体
type SystemInfo struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Uptime      time.Duration `json:"uptime"`
	Memory      string        `json:"memory"`
	CPU         string        `json:"cpu"`
	Connections int           `json:"connections"`
	Channels    int           `json:"channels"`
	Devices     int           `json:"devices"`
}

var startTime = time.Now()

// GetUptime 获取服务的运行时长
func GetUptime() time.Duration {
	return time.Since(startTime)
}

func SystemRoutes(rg *gin.RouterGroup) {
	rg.GET("/system/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, Resp(GetSystemInfo()))
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() SystemInfo {
	var info SystemInfo

	info.Name = "wiTeam数采网关"

	// 获取 Go 版本
	info.Version = runtime.Version()

	// 获取系统运行时间
	info.Uptime = GetUptime()

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.Memory = fmt.Sprintf("Total: %v MB, Free: %v MB", memInfo.Total/(1024*1024), memInfo.Free/(1024*1024))
	}

	// 获取 CPU 信息
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil {
		info.CPU = fmt.Sprintf("%.2f%%", cpuPercent[0])
	}

	// 获取网络连接数
	info.Connections = 0

	// 这里 Channels 暂时模拟为 0，实际应用中需要根据具体逻辑获取
	info.Channels = 0

	// 获取网络设备信息

	info.Devices = 0
	return info
}
