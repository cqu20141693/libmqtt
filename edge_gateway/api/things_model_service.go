package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/edge_gateway/orm"
	"net/http"
	"strconv"
)

func ThingsModelRoutes(rg *gin.RouterGroup) {
	// 查询设备下物模型
	deviceThingsModel(rg)

	// 新增物模型
	addThingsModel(rg)
	// 删除物模型
	deleteThingsModel(rg)
	// 更新物模型
	updateThignsModel(rg)

}
func updateThignsModel(rg *gin.RouterGroup) gin.IRoutes {
	return rg.PUT("/things-model/:id", func(c *gin.Context) {
		var info orm.DeviceThingsModel

		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
		}
		if info.ID == 0 {
			info.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
		}
		err = orm.DB.ThingsModelUpdate(info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据更新失败"))
		}
		c.JSON(http.StatusOK, Resp("更新成功"))
	})
}
func deleteThingsModel(rg *gin.RouterGroup) gin.IRoutes {
	return rg.DELETE("/things-model/:id", func(c *gin.Context) {

		ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("id 不存在"))
			return
		}

		err = orm.DB.ThingsModelDelete(ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据删除失败"))
			return
		}
		c.JSON(http.StatusOK, Resp("数据删除成功"))
	})
}

func addThingsModel(rg *gin.RouterGroup) gin.IRoutes {
	return rg.POST("/things-model", func(c *gin.Context) {
		var info orm.DeviceThingsModel
		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
			return
		}
		var id int64
		id, err = orm.DB.ThingsModelAdd(info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("新增失败"))
			return
		}
		info.ID = id
		c.JSON(http.StatusOK, Resp("新增成功"))
	})
}

func deviceThingsModel(rg *gin.RouterGroup) gin.IRoutes {
	return rg.GET("/things/list/:deviceId", func(c *gin.Context) {
		deviceId, err := strconv.ParseInt(c.Param("deviceId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
			return
		}
		thingsModels, err := orm.DB.ThingsModelByDevice(deviceId)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("查询失败"))
			return
		}
		c.JSON(http.StatusOK, Resp(thingsModels))
	})
}
