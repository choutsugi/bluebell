package test

import (
	"bluebell/dao/redis"
	"testing"
)

func TestGoRedis(t *testing.T) {
	if err := redis.InitRDB(); err != nil {
		t.Error(err)
	}
}
