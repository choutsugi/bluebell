package main

import (
	"bluebell/internal/conf"
	"bluebell/internal/server"
	"flag"
)

var flagConf string

func init() {
	flag.StringVar(&flagConf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {

	flag.Parse()

	if err := conf.Load(flagConf); err != nil {
		panic(err)
	}

	app := server.NewServer(conf.Boot)
	app.Run()
}
