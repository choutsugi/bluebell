package conf

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Bootstrap struct {
	App       *App       `mapstructure:"app"`
	Log       *Log       `mapstructure:"log"`
	Data      *Data      `mapstructure:"data"`
	SnowFlake *SnowFlake `mapstructure:"snowflake"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type Data struct {
	Cache      *Cache      `mapstructure:"cache"`
	DataSource *DataSource `mapstructure:"datasource"`
}

type Cache struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type DataSource struct {
	Dsn string `mapstructure:"dsn"`
}

type SnowFlake struct {
	StartTime string `mapstructure:"start_time"`
	MachineId int64  `mapstructure:"machine_id"`
}

var Boot = new(Bootstrap)

func Load() (err error) {
	filePath := flag.String("conf", "../configs/config.yaml", "config path, eg: -conf config.yaml")

	flag.Parse()
	viper.SetConfigFile(*filePath)

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(Boot); err != nil {
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Boot); err != nil {
			log.Printf("viper unmarshal config-file failed, err:%v\n", err)
		}
	})

	return
}
