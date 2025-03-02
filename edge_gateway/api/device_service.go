package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/edge_gateway/orm"
	"net/http"
	"strconv"
)

func DeviceRoutes(rg *gin.RouterGroup) {
	// 查询通道下设备
	channelDevice(rg)

	// 新增设备
	addDevice(rg)
	// 删除设备
	deleteDevice(rg)
	// 更新设备
	updateDevice(rg)

}
func updateDevice(rg *gin.RouterGroup) gin.IRoutes {
	return rg.PUT("/device/:id", func(c *gin.Context) {
		var info orm.Device

		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
		}
		if info.ID == 0 {
			info.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
		}
		err = orm.DB.DeviceUpdate(info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据更新失败"))
		}
		c.JSON(http.StatusOK, Resp("更新成功"))
	})
}
func deleteDevice(rg *gin.RouterGroup) gin.IRoutes {
	return rg.DELETE("/device/:id", func(c *gin.Context) {

		ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("id 不存在"))
			return
		}

		err = orm.DB.DeviceDelete(ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据删除失败"))
			return
		}
		c.JSON(http.StatusOK, Resp("数据删除成功"))
	})
}

func addDevice(rg *gin.RouterGroup) gin.IRoutes {
	return rg.POST("/device", func(c *gin.Context) {
		var info orm.Device
		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
			return
		}
		var id int64
		id, err = orm.DB.DeviceAdd(info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("新增失败"))
			return
		}
		info.ID = id
		c.JSON(http.StatusOK, Resp("新增成功"))
	})
}

func channelDevice(rg *gin.RouterGroup) gin.IRoutes {
	return rg.GET("/device/list/:channelId", func(c *gin.Context) {
		channelId, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
			return
		}
		devices, err := orm.DB.DeviceByChannel(channelId)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("查询失败"))
			return
		}
		c.JSON(http.StatusOK, Resp(devices))
	})
}
