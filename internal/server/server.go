package server

import (
	"bluebell/internal/conf"
	"bluebell/internal/data"
	"bluebell/internal/router"
	"bluebell/pkg/logger"
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
	config := logger.Config{
		Level:      c.Log.Level,
		FileName:   c.Log.FileName,
		MaxSize:    c.Log.MaxSize,
		MaxAge:     c.Log.MaxAge,
		MaxBackups: c.Log.MaxBackups,
	}
	if err := logger.Init(config); err != nil {
		panic(err)
	}

	//建立数据库连接
	db := data.NewDataSource(c.Data.DataSource)
	rdb := data.NewCache(c.Data.Cache)
	_ = data.NewData(db, rdb)

	srv := &Server{
		addr:    c.App.Addr,
		handler: router.Setup(),
	}
	return srv
}
