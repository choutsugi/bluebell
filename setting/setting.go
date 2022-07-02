package setting

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App   *AppConfig   `mapstructure:"app"`
	Log   *LogConfig   `mapstructure:"log"`
	Db    *DbConfig    `mapstructure:"db"`
	Redis *RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type DbConfig struct {
	DriveName   string `mapstructure:"drive_name"`
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"db_name"`
	Port        int    `mapstructure:"port"`
	MaxOpenCons int    `mapstructure:"max_open_cons"`
	MaxIdleCons int    `mapstructure:"max_idle_cons"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

var Conf = new(Config)

func Init() (err error) {
	filePath := flag.String("conf", "./config/config.yaml", "config path, eg: -conf config.yaml")

	flag.Parse()
	viper.SetConfigFile(*filePath)

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(Conf); err != nil {
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper unmarshal config-file failed, err:%v\n", err)
		}
	})

	return
}
