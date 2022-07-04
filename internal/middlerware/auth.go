package middlerware

import (
	"bluebell/internal/pkg/auth"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Request.Header.Get("Authorization")
		if tokenStr == "" {
			result.Error(ctx, errx.ErrTokenMissing)
			ctx.Abort()
			return
		}
		tokenStr = tokenStr[len(auth.Conf.TokenType)+1:]
		// Token 解析校验
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

		claims := token.Claims.(*auth.Claims)
		if claims.Issuer != auth.Conf.Issuer {
			result.Error(ctx, errx.ErrTokenInvalid)
			ctx.Abort()
			return
		}

		uid, err := strconv.ParseInt(claims.ID, 10, 64)
		if err != nil {
			result.Error(ctx, errx.ErrTokenInvalid)
			ctx.Abort()
			return
		}
		ctx.Set("uid", uid)
	}
}
