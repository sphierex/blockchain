package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers manages the set of bar ledger endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Sample just provides a starting point for the class.
func (h Handlers) Sample(ctx *gin.Context) {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	ctx.JSON(http.StatusOK, resp)
}
