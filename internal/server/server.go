package server

import (
	v1 "bluebell/api/v1"
	"bluebell/internal/conf"
	"bluebell/internal/data"
	"bluebell/internal/pkg/auth"
	"bluebell/internal/repository"
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

func NewServer(c *conf.Bootstrap) *Server {

	//初始化日志
	logConfig := logger.Config{
		Level:      c.Log.Level,
		FileName:   c.Log.FileName,
		MaxSize:    c.Log.MaxSize,
		MaxAge:     c.Log.MaxAge,
		MaxBackups: c.Log.MaxBackups,
	}
	if err := logger.Init(logConfig); err != nil {
		panic(err)
	}

	//初始化雪花算法
	if err := snowflake.Init(c.SnowFlake.StartTime, c.SnowFlake.MachineId); err != nil {
		panic(err)
	}

	//建立数据库连接
	dataSource := data.NewDataSource(c.Data.DataSource)
	cache := data.NewCache(c.Data.Cache)
	database := data.NewData(dataSource, cache)

	//初始化jwt
	jwtConfig := auth.Config{
		TokenType:            c.Jwt.TokenType,
		Issuer:               c.Jwt.Issuer,
		Secret:               c.Jwt.Secret,
		TTL:                  c.Jwt.TTL,
		BlacklistKeyPrefix:   c.Jwt.BlacklistKeyPrefix,
		BlacklistGracePeriod: c.Jwt.BlacklistGracePeriod,
		RefreshGracePeriod:   c.Jwt.RefreshGracePeriod,
		RefreshLockName:      c.Jwt.RefreshLockName,
	}
	auth.Init(jwtConfig, cache)

	userRepo := repository.NewUserRepo(database.DB)
	userService := service.NewUserService(userRepo)

	api := v1.Register(userService)

	srv := &Server{
		addr:    c.App.Addr,
		handler: router.Setup(api),
	}
	return srv
}
