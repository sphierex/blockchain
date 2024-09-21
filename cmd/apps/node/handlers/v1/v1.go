package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sphierex/blockchain/cmd/apps/node/handlers/v1/private"
	"github.com/sphierex/blockchain/cmd/apps/node/handlers/v1/public"
	"go.uber.org/zap"
)

const version = "v1"

type Config struct {
	Log *zap.SugaredLogger
}

// PublicRoutes binds all the version 1 public routes.
func PublicRoutes(app *gin.Engine, cfg Config) {
	pbl := public.Handlers{Log: cfg.Log}

	v1 := app.Group(version)
	{
		v1.GET("/sample", pbl.Sample)
	}
}

// PrivateRoutes binds all the version 1 private routes.
func PrivateRoutes(app *gin.Engine, cfg Config) {
	prv := private.Handlers{
		Log: cfg.Log,
	}

	v1 := app.Group(version)
	{
		v1.GET("/node/sample", prv.Sample)
	}
}
