package mid

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Infow("request started", "traceId", ctx.GetString("tradeId"), "method", ctx.Request.Method,
			"path", ctx.Request.URL.Path, "remoteAddr", ctx.RemoteIP())

		ctx.Next()

		log.Infow("request started", "traceId", ctx.GetString("tradeId"), "method", ctx.Request.Method,
			"path", ctx.Request.URL.Path, "remoteAddr", ctx.RemoteIP())
	}
}
