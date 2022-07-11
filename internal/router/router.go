package router

import (
	"bluebell/api/v1"
	"bluebell/internal/conf"
	"bluebell/internal/middlerware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Setup(api v1.Api, config *conf.Bootstrap) *gin.Engine {
	r := gin.New()

	r.Use(
		middlerware.Logger(),
		middlerware.Recovery(true),
		middlerware.RateLimit(time.Second*time.Duration(config.RateLimit.FillInterval), config.RateLimit.Capacity),
	)
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

		certified.GET("community/all", api.Community.FetchAll)
		certified.GET("community/:id", api.Community.FetchOneById)

		certified.GET("post/all", api.Post.FetchAll)
		certified.GET("post/list/order", api.Post.FetchListWithOrder)
		certified.GET("post/list/paginate", api.Post.FetchListByPaginate)
		certified.GET("post/:id", api.Post.FetchById)
		certified.POST("post", api.Post.Create)
		certified.PUT("post", api.Post.Update)
		certified.DELETE("post/:id", api.Post.Delete)

		certified.POST("vote", api.Vote.PostVote)
	}

	return r
}
