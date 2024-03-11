package bootstrap

import pkg "orchestrator-service/pkg/logger"

func GetLogger() pkg.Logger {
	return pkg.NewLogrus()
}
