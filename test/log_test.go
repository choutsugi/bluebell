package test

import (
	"bluebell/internal/conf"
	"bluebell/pkg/logger"
	"go.uber.org/zap"
	"testing"
)

func TestLog(t *testing.T) {
	if err := conf.Load(); err != nil {
		t.Error(err)
	}

	config := logger.Config{
		Level:      conf.Boot.Log.Level,
		FileName:   conf.Boot.Log.FileName,
		MaxSize:    conf.Boot.Log.MaxSize,
		MaxAge:     conf.Boot.Log.MaxAge,
		MaxBackups: conf.Boot.Log.MaxBackups,
	}
	if err := logger.Init(config); err != nil {
		t.Error(err)
	}
	defer logger.Sync()

	zap.L().Info("just test for logger")
}
