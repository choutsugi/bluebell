package test

import (
	"bluebell/internal/conf"
	"bluebell/internal/data"
	"testing"
)

func TestCache(t *testing.T) {
	if err := conf.Load(); err != nil {
		t.Error(err)
	}

	data.NewCache(conf.Boot.Data.Cache)
}

func TestDataSource(t *testing.T) {
	if err := conf.Load(); err != nil {
		t.Error(err)
	}

	data.NewDataSource(conf.Boot.Data.DataSource)
}
