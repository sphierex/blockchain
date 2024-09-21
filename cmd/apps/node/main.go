package main

import (
	"go.uber.org/zap"

	"github.com/sphierex/blockchain/pkg/logger"
)

var build = "develop"

func main() {
	log, err := logger.New("NODE")
	if err != nil {
		return
	}
	defer func(log *zap.SugaredLogger) {
		_ = log.Sync()
	}(log)

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
	}
}

func run(log *zap.SugaredLogger) error {
	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	return nil
}
