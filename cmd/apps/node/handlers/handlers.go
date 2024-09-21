package handlers

import (
	"expvar"
	"github.com/sphierex/blockchain/internal/web/mid"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sphierex/blockchain/cmd/apps/node/handlers/debug/checkgrp"
	v1 "github.com/sphierex/blockchain/cmd/apps/node/handlers/v1"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// MuxConfig contains all the mandatory systems required by handlers.
type MuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// PublicMux constructs a http.Handler with all application routes defined.
func PublicMux(cfg MuxConfig) http.Handler {
	app := gin.New()
	app.Use(
		mid.TraceId(),
		mid.Logger(cfg.Log),
		mid.Metrics(),
		gin.Recovery(),
	)

	v1.PublicRoutes(app, v1.Config{Log: cfg.Log})

	return app
}

// PrivateMux constructs a http.Handler with all application routes defined.
func PrivateMux(cfg MuxConfig) http.Handler {
	app := gin.New()
	app.Use(
		mid.TraceId(),
		mid.Logger(cfg.Log),
		mid.Metrics(),
		gin.Recovery(),
	)

	v1.PrivateRoutes(app, v1.Config{Log: cfg.Log})

	return app
}

// standardDebugMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func standardDebugMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := standardDebugMux()

	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}
