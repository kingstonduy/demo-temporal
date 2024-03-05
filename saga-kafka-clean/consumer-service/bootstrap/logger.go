package bootstrap

import "consumer-service/pkg/logger"

func GetLogger() logger.Logger {
	return logger.NewLogrus()
}
