package middlerware

import (
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

// RateLimit 限流中间件
func RateLimit(fillInterval time.Duration, cap int64) func(ctx *gin.Context) {

	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) > 0 {
			result.Error(ctx, errx.ErrInternalServerBusy)
			return
		}
		ctx.Next()
	}
}
