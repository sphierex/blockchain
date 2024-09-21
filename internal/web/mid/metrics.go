package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/sphierex/blockchain/internal/web/metrics"
)

func Metrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := metrics.Set(ctx)

		ctx.Next()

		metrics.AddRequests(c)
		metrics.AddGoroutines(c)

		if len(ctx.Errors) > 0 {
			metrics.AddErrors(c)
		}
	}

}
