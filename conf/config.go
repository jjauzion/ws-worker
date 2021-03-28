package conf

import (
	"github.com/jjauzion/ws-worker/internal/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	WS_GRPC_HOST            string
	WS_GRPC_PORT            string
	WS_DOCKER_LOG_FOLDER    string
	WS_DOCKER_RESULT_FOLDER string
}

func GetConfig(log *logger.Logger) (Configuration, error) {
	cf := Configuration{}
	err := viper.Unmarshal(&cf)
	if err != nil {
		return cf, err
	}
	log.Info("configuration loaded")
	return cf, nil
}
