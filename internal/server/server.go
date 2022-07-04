package server

import (
	"bluebell/internal/conf"
	"bluebell/router"
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
	srv := &Server{
		addr:    c.App.Addr,
		handler: router.InitRouter(),
	}
	return srv
}
