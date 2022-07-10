package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Bootstrap struct {
	App       *App       `mapstructure:"app"`
	Log       *Log       `mapstructure:"log"`
	Data      *Data      `mapstructure:"data"`
	SnowFlake *SnowFlake `mapstructure:"snowflake"`
	Jwt       *Jwt       `mapstructure:"jwt"`
	Ranking   *Ranking   `mapstructure:"ranking"`
}

type App struct {
	Name string `mapstructure:"name"`
	Addr string `mapstructure:"addr"`
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
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	Db           int           `mapstructure:"db"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DataSource struct {
	Dsn string `mapstructure:"dsn"`
}

type SnowFlake struct {
	StartTime string `mapstructure:"start_time"`
	MachineId int64  `mapstructure:"machine_id"`
}

type Jwt struct {
	TokenType            string `mapstructure:"token_type"`
	Issuer               string `mapstructure:"issuer"`
	Secret               string `mapstructure:"secret"`
	TTL                  int64  `mapstructure:"ttl"`                    // 有效时间
	BlacklistKeyPrefix   string `mapstructure:"blacklist_key_prefix"`   // 黑名单Key前缀
	BlacklistGracePeriod int64  `mapstructure:"blacklist_grace_period"` // 黑名单宽限时间（避免并发请求失败）
	RefreshGracePeriod   int64  `mapstructure:"refresh_grace_period"`
	RefreshLockName      string `mapstructure:"refresh_lock_name"`
}

type Ranking struct {
	PostTimeKey            string  `mapstructure:"post_time_key"`
	PostScoreKey           string  `mapstructure:"post_score_key"`
	PostVotedKeyPrefix     string  `mapstructure:"post_voted_key_prefix"`
	PostCommunityKeyPrefix string  `mapstructure:"post_community_key_prefix"`
	PostVotingPeriod       float64 `mapstructure:"post_voting_period"`
	PostVoteUnitScore      float64 `mapstructure:"post_vote_unit_score"`
}

var Boot = new(Bootstrap)

func Load(flagConf string) (err error) {

	viper.SetConfigFile(flagConf)

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
