package test

import "flag"

var flagConf string

func init() {
	flag.StringVar(&flagConf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}
