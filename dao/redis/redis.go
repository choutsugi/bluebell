package redis

import "github.com/go-redis/redis"

var RDB *redis.Client

func InitClient() (err error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "dangerous",
		DB:       0,
		PoolSize: 100,
	})

	_, err = RDB.Ping().Result()
	return
}
