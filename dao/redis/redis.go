package redis

import (
	"bluebell/setting"
	"fmt"
	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRDB(config *setting.RedisConfig) (err error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
		PoolSize: config.PoolSize,
	})

	_, err = RDB.Ping().Result()
	return
}

func Close() {
	RDB.Close()
}
