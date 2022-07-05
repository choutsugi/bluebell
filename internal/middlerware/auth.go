package middlerware

import (
	"bluebell/internal/pkg/auth"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//解析token
		tokenStr := ctx.Request.Header.Get("Authorization")
		if tokenStr == "" {
			result.Error(ctx, errx.ErrTokenMissing)
			ctx.Abort()
			return
		}
		tokenStr = tokenStr[len(auth.Conf.TokenType)+1:]
		token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(auth.Conf.Secret), nil
		})
		if err != nil {
			//if ve, ok := err.(*jwt.ValidationError); ok {
			//	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			//	} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
			//	} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
			//	} else {
			//	}
			//}
			result.Error(ctx, errx.ErrTokenInvalid)
			ctx.Abort()
			return
		}

		//校验发行人
		claims := token.Claims.(*auth.Claims)
		if claims.Issuer != auth.Conf.Issuer {
			result.Error(ctx, errx.ErrTokenInvalid)
			ctx.Abort()
			return
		}

		//黑名单过滤
		if auth.IsInBlacklist(tokenStr) {
			result.Error(ctx, errx.ErrTokenInvalid)
			ctx.Abort()
			return
		}

		//续签
		if claims.ExpiresAt.Unix()-time.Now().Unix() < auth.Conf.RefreshGracePeriod {
			//TODO：续签之前需要检查用户是否仍存在于数据库中
			tokenData, err := auth.GenerateToken(claims.UID)
			if err != nil {
				result.Error(ctx, errx.ErrInternalServerError)
				ctx.Abort()
				return
			}
			ctx.Header("new_access_token", tokenData.AccessToken)
			ctx.Header("new_expired_in", strconv.FormatInt(tokenData.ExpiresIn, 10))
			ctx.Header("new_token_type", tokenData.TokenType)
		}

		ctx.Set("uid", claims.UID)
	}
}
