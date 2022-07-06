package data

import (
	"bluebell/internal/conf"
	"bluebell/internal/entity"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Data struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewData(db *gorm.DB, rdb *redis.Client) *Data {
	return &Data{
		DB:  db,
		RDB: rdb,
	}
}

func NewDataSource(c *conf.DataSource) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(c.Dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Panicf("failed to establish connection with database: %v", err)
	}

	if err = db.AutoMigrate(&entity.User{}, &entity.Community{}, &entity.Post{}); err != nil {
		log.Panicf("failed to migrate table to database: %v", err)
	}
	return db
}

func NewCache(c *conf.Cache) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.Db,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	})

	status := rdb.Ping()
	if err := status.Err(); err != nil {
		log.Panicf("failed to establish connection with redis: %v", err)
	}
	return rdb
}
