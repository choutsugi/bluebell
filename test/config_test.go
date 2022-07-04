package test

import (
	"bluebell/internal/conf"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {

	if err := conf.Load(); err != nil {
		t.Error(err)
	}

	fmt.Println(*conf.Boot.App)
	fmt.Println(*conf.Boot.Log)
	fmt.Println(*conf.Boot.Data.Cache)
	fmt.Println(*conf.Boot.Data.DataSource)
	fmt.Println(*conf.Boot.SnowFlake)

}
