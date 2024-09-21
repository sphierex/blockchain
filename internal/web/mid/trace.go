package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := uuid.New().String()
		ctx.Set("tradeId", traceId)
		ctx.Next()
	}
}
