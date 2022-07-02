package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/setting"
)

func main() {
	if err := setting.Init(); err != nil {
		panic(err)
	}

	if err := mysql.InitDB(setting.Conf.Db); err != nil {
		panic(err)
	}

	if err := redis.InitRDB(setting.Conf.Redis); err != nil {
		panic(err)
	}

}
