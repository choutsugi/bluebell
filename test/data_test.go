package test

import (
	"bluebell/internal/conf"
	"bluebell/internal/data"
	"flag"
	"testing"
)

func TestCache(t *testing.T) {
	flag.Parse()
	if err := conf.Load(flagConf); err != nil {
		t.Error(err)
	}

	data.NewCache(conf.Boot.Data.Cache)
}

func TestDataSource(t *testing.T) {
	flag.Parse()
	if err := conf.Load(flagConf); err != nil {
		t.Error(err)
	}

	data.NewDataSource(conf.Boot.Data.DataSource)
}
