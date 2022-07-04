package router

import (
	"bluebell/internal/middlerware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(middlerware.Logger(), middlerware.Recovery(true))

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Page Not Found",
		})
	})

	return r
}
