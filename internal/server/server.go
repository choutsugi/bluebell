package server

import (
	v1 "bluebell/api/v1"
	"bluebell/internal/conf"
	"bluebell/internal/data"
	"bluebell/internal/data/cache"
	"bluebell/internal/data/repo"
	"bluebell/internal/pkg/auth"
	"bluebell/internal/router"
	"bluebell/internal/service"
	"bluebell/pkg/logger"
	"bluebell/pkg/snowflake"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	addr    string
	handler *gin.Engine
}

func (srv *Server) Run() {

	defer logger.Sync()

	s := http.Server{
		Addr:    srv.addr,
		Handler: srv.handler,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("server stopped unexpectedly, err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := s.Shutdown(ctx); err != nil {
		log.Panicf("server stopped unexpectedly, err: %v", err)
	}
}

func NewServer(config *conf.Bootstrap) *Server {

	//初始化日志
	logConfig := logger.Config{
		Level:      config.Log.Level,
		FileName:   config.Log.FileName,
		MaxSize:    config.Log.MaxSize,
		MaxAge:     config.Log.MaxAge,
		MaxBackups: config.Log.MaxBackups,
	}
	if err := logger.Init(logConfig); err != nil {
		panic(err)
	}

	//初始化雪花算法
	if err := snowflake.Init(config.SnowFlake.StartTime, config.SnowFlake.MachineId); err != nil {
		panic(err)
	}

	//建立数据库连接
	db := data.NewDataSource(config.Data.DataSource)
	rdb := data.NewCache(config.Data.Cache)
	database := data.NewData(db, rdb)

	//初始化jwt
	jwtConfig := auth.Config{
		TokenType:            config.Jwt.TokenType,
		Issuer:               config.Jwt.Issuer,
		Secret:               config.Jwt.Secret,
		TTL:                  config.Jwt.TTL,
		BlacklistKeyPrefix:   config.Jwt.BlacklistKeyPrefix,
		BlacklistGracePeriod: config.Jwt.BlacklistGracePeriod,
		RefreshGracePeriod:   config.Jwt.RefreshGracePeriod,
		RefreshLockName:      config.Jwt.RefreshLockName,
	}
	auth.Init(jwtConfig, rdb)

	//Repo
	userRepo := repo.NewUserRepo(database.DB)
	communityRepo := repo.NewCommunityRepo(database.DB)
	postRepo := repo.NewPostRepo(database.DB)

	//Cache
	voteCache := cache.NewVoteCache(rdb, config.Ranking)

	//Service
	userService := service.NewUserService(userRepo)
	communityService := service.NewCommunityService(communityRepo)
	postService := service.NewPostService(postRepo, userRepo, communityRepo, voteCache)
	voteService := service.NewVoteService(voteCache)

	//Register Services
	api := v1.Register(userService, communityService, postService, voteService)

	srv := &Server{
		addr:    config.App.Addr,
		handler: router.Setup(api, config),
	}
	return srv
}
