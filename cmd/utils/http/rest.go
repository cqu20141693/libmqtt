package http

import "github.com/gin-gonic/gin"

func ResponseSuccess(g *gin.Context, data any) {
	g.JSON(200, gin.H{
		"code":   "0000",
		"result": &data,
		"msg":    "success",
	})
}

func Response400(g *gin.Context, code string, msg string) {
	failResponse(g, 400, code, msg)
}

func Response401(g *gin.Context, code string, msg string) {
	failResponse(g, 401, code, msg)
}
func Response403(g *gin.Context, code string, msg string) {
	failResponse(g, 403, code, msg)
}

func Response500(g *gin.Context, code string, msg string) {
	failResponse(g, 500, code, msg)
}

func failResponse(g *gin.Context, c int, code string, msg string) {
	g.JSON(c, gin.H{
		"code": code,
		"msg":  msg,
	})
}
