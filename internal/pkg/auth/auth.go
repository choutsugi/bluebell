package auth

import (
	"bluebell/internal/pkg/errx"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type Config struct {
	TokenType            string
	Issuer               string
	Secret               string
	TTL                  int64
	BlacklistKeyPrefix   string
	BlacklistGracePeriod int64
	RefreshGracePeriod   int64
	RefreshLockName      string
}

type TokenData struct {
	TokenType   string
	AccessToken string
	ExpiresIn   int64
}

type CustomClaims struct {
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

	accessClaim := CustomClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Conf.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(Conf.TTL) * time.Second)),
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
	tokenData.ExpiresIn = Conf.TTL
	return
}

func IsInBlacklist(token string) bool {
	joinTimeStr, err := cache.Get(getBlacklistKey(token)).Result()
	if err != nil {
		return false
	}
	joinTime, err := strconv.ParseInt(joinTimeStr, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix()-joinTime < Conf.BlacklistGracePeriod {
		return false
	}

	return true
}

func JoinBlacklist(token *jwt.Token) error {
	claim, ok := token.Claims.(*CustomClaims)
	if !ok {
		return errx.ErrTokenInvalid
	}

	timer := claim.ExpiresAt.Unix() - time.Now().Unix()
	return cache.SetNX(getBlacklistKey(token.Raw), time.Now().Unix(), time.Duration(timer)*time.Second).Err()
}

func getBlacklistKey(token string) string {
	return Conf.BlacklistKeyPrefix + token
}
