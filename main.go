package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/setting"
	"fmt"
	"go.uber.org/zap"
)

func main() {

	var err error

	if err = setting.Init(); err != nil {
		panic(err)
	}

	if err = mysql.InitDB(setting.Conf.Db); err != nil {
		panic(err)
	}

	if err = redis.InitRDB(setting.Conf.Redis); err != nil {
		panic(err)
	}

	if err = logger.InitLogger(setting.Conf.Log); err != nil {
		panic(err)
	}

	if err = snowflake.Init(setting.Conf.SnowFlake.StartTime, setting.Conf.SnowFlake.MachineId); err != nil {
		zap.L().Panic("Init snowflake error", zap.Error(err))
	}

	fmt.Println(snowflake.GenerateID())
}
