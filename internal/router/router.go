package router

import (
	v1 "bluebell/api/v1"
	"bluebell/internal/middlerware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup(api v1.Api) *gin.Engine {
	r := gin.New()

	r.Use(middlerware.Logger(), middlerware.Recovery(true))

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Page Not Found",
		})
	})

	r.POST("/api/v1/user/signup", api.User.Signup)
	r.POST("/api/v1/user/login", api.User.Login)

	admin := r.Group("/api/v1/auth/", middlerware.JwtAuth())
	{
		admin.DELETE("logout", api.User.Logout)
		admin.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	return r
}
