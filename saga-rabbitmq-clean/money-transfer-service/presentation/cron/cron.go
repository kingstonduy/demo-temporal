package cron

import (
	"context"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

const (
	TopicAccountStatement = "OCBDWDAILY.TODAY_ACCOUNT_STATEMENT"
)

type CronServer struct {
	logger logger.Logger
	broker broker.Broker
}

func NewCronServer(logger logger.Logger, broker broker.Broker) *CronServer {
	return &CronServer{
		logger: logger,
		broker: broker,
	}
}

func (s *CronServer) Start(ctx context.Context) error {
	opts := s.GetStartOptions()
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func(ctx context.Context) {
			s.logger.Info(ctx, "Stop Cron Server")
		}(ctx)
		<-ctx.Done()
	}()

	s.logger.Infof(ctx, "Started Cron Server")
	return nil
}
