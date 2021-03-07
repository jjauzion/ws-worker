package client

import (
	"github.com/jjauzion/ws-worker/conf"
	"github.com/jjauzion/ws-worker/internal/logger"
	"log"
)

func dependencies() (*logger.Logger, conf.Configuration, error) {
	lg, err := logger.ProvideLogger()
	if err != nil {
		log.Fatalf("cannot create logger %v", err)
	}

	cf, err := conf.GetConfig(lg)
	if err != nil {
		log.Fatalf("cannot get config %v", err)
	}

	return lg, cf, nil
}
