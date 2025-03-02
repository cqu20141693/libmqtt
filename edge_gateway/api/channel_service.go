package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/edge_gateway/constants"
	"github.com/goiiot/libmqtt/edge_gateway/orm"
	"net/http"
	"strconv"
)

func ChannelRoutes(rg *gin.RouterGroup) {
	// 通道树
	channelTree(rg)

	// 通道组添加
	addChannelGroup(rg)

	// 添加通道
	addChannel(rg)
	updateChannel(rg)
	rg.DELETE("/channel/:id", func(c *gin.Context) {

		ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("id 不存在"))
			return
		}

		err = orm.DB.ChannelDelete(ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据删除失败"))
			return
		}
		c.JSON(http.StatusOK, Resp("数据删除成功"))
	})

}

func updateChannel(rg *gin.RouterGroup) gin.IRoutes {
	return rg.PUT("/channel/:id", func(c *gin.Context) {
		var info orm.Channel
		info.ParentID = DefaultParentId

		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
		}
		if info.ID == 0 {
			info.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
		}
		err = orm.DB.ChannelUpdate(info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据更新失败"))
		}
		c.JSON(http.StatusOK, Resp("更新成功"))
	})
}

func addChannel(rg *gin.RouterGroup) gin.IRoutes {
	return rg.POST("/channel", func(c *gin.Context) {
		var info orm.Channel
		info.ParentID = DefaultParentId
		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
			return
		}
		addChannel, err := orm.DB.ChannelAdd(info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据保存失败"))
			return
		}
		info.ID = addChannel
		c.JSON(http.StatusOK, Resp(info))
	})
}

func addChannelGroup(rg *gin.RouterGroup) gin.IRoutes {
	return rg.POST("/channel/group", func(c *gin.Context) {
		type ChannelGroup struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description" `
			ParentId    int64  `json:"parentId"`
		}
		var info ChannelGroup
		err := c.BindJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, RespError("参数错误"))
		}

		pId := DefaultParentId
		if info.ParentId > 0 {
			pId = info.ParentId
		}

		channel := orm.Channel{
			Name:     info.Name,
			Type:     constants.ChannelGroupKey,
			ParentID: pId,
			Config:   "",
			Desc:     info.Description,
		}
		addChannel, err := orm.DB.ChannelAdd(channel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据保存失败"))
			return
		}
		channel.ID = addChannel
		c.JSON(http.StatusOK, Resp(channel))
	})
}

func channelTree(rg *gin.RouterGroup) gin.IRoutes {
	return rg.GET("/channel/tree", func(c *gin.Context) {
		channels, err := orm.DB.ChannelAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, RespError("数据查询失败"))
			return
		}
		c.JSON(http.StatusOK, Resp(channels))
	})
}
