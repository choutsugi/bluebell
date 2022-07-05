package auth

import (
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Config struct {
	TokenType            string
	Issuer               string
	Secret               string
	TTL                  time.Duration
	BlacklistKeyPrefix   string
	BlacklistGracePeriod time.Duration
	RefreshGracePeriod   int64
	RefreshLockName      string
}

type TokenData struct {
	TokenType   string
	AccessToken string
	ExpiresIn   int64
}

type Claims struct {
	UID int64
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

func GenerateToken(uid int64) (tokenData TokenData, err error) {

	accessClaim := Claims{
		uid,
		jwt.RegisteredClaims{
			Issuer:    Conf.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Conf.TTL)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tokenData.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim).
		SignedString([]byte(Conf.Secret))
	if err != nil {
		return TokenData{}, err
	}
	tokenData.TokenType = Conf.TokenType
	tokenData.ExpiresIn = int64(Conf.TTL / time.Second)

	return
}

func IsInBlacklist(token string) bool {
	if _, err := cache.Get(getBlacklistKey(token)).Result(); err != nil {
		return false
	}
	return true
}

func JoinBlacklist(token *jwt.Token) error {
	curUnixTime := time.Now().Unix()
	timer := time.Duration(token.Claims.(Claims).ExpiresAt.Unix()-curUnixTime) * time.Second
	return cache.SetNX(getBlacklistKey(token.Raw), curUnixTime, timer).Err()
}

func getBlacklistKey(token string) string {
	return Conf.BlacklistKeyPrefix + token
}
