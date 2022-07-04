package test

import (
	"bluebell/internal/conf"
	"flag"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {

	flag.Parse()
	if err := conf.Load(flagConf); err != nil {
		t.Error(err)
	}

	fmt.Println(conf.Boot)

}
