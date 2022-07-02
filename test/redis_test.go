package test

import (
	"bluebell/dao/redis"
	"testing"
)

func TestGoRedis(t *testing.T) {
	if err := redis.InitClient(); err != nil {
		t.Error(err)
	}
}
