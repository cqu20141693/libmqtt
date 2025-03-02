package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"net/http"
)

func UserRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/login", func(c *gin.Context) {

		type UserLogin struct {
			Username string `form:"username" json:"username" xml:"username" binding:"required"`
			Password string `form:"password" json:"password" xml:"password" binding:"required"`
		}
		var login UserLogin
		err := c.BindJSON(&login)
		if err != nil {
			cclog.SugarLogger.Error(err)
			c.JSON(http.StatusBadRequest, "请求参数错误")
			return
		}

		data := make(map[string]interface{})
		data["token"] = "go-gateway"
		c.JSON(http.StatusOK, Resp(data))
	})
}
