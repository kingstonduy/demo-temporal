package bootstrap

import "orchestrator-service/pkg/logger"

func GetLogger() logger.Logger {
	return logger.NewLogrus()
}
