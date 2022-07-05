package router

import (
	"bluebell/api/v1"
	"bluebell/internal/middlerware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup(api v1.Api) *gin.Engine {
	r := gin.New()

	r.Use(middlerware.Logger(), middlerware.Recovery(true))
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "404",
		})
	})

	group := r.Group("/api/v1")
	{
		group.POST("signup", api.User.Signup)
		group.POST("login", api.User.Login)
	}

	certified := group.Group("auth", middlerware.JwtAuth())
	{
		certified.DELETE("logout", api.User.Logout)
		certified.GET("communities", api.Community.FetchAll)
		certified.GET("community/:id", api.Community.FetchOneById)
	}

	return r
}
