package bootstrap

import "orchestrator-service/infra/logger"

func GetLogger() logger.Logger {
	return logger.NewLogrus()
}
