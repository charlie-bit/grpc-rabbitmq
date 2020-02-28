package controller

import "github.com/gin-gonic/gin"

func SayHello(json string) string {
	return "hello " + json
}

func Hello(ctx *gin.Context) {
	ctx.JSON(200,"hello")
}
