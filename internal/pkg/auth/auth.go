package auth

import (
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type Config struct {
	TokenType            string
	Issuer               string
	Secret               string
	TTL                  time.Duration
	BlacklistKeyPrefix   string
	BlacklistGracePeriod time.Duration
	RefreshGracePeriod   time.Duration
	RefreshLockName      string
}

type Result struct {
	TokenType   string
	AccessToken string
	ExpiresIn   int64
}

type Claims struct {
	jwt.RegisteredClaims
}

var (
	Conf  Config
	cache *redis.Client
)

func Init(config Config, rdb *redis.Client) {
	Conf = config
	cache = rdb
}

func GenerateToken(uid int64) (ret Result, err error) {

	claim := Claims{
		jwt.RegisteredClaims{
			Issuer:    Conf.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Conf.TTL)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strconv.FormatInt(uid, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ret.AccessToken, err = token.SignedString([]byte(Conf.Secret))
	ret.TokenType = Conf.TokenType
	ret.ExpiresIn = int64(Conf.TTL / time.Second)
	return
}
