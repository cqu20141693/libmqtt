package gaea

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddPingRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/ping")

	ping.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}
