package broker

import (
	"context"
	"encoding/json"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/pipeline"
	"github.com/lengocson131002/go-clean-core/transport/broker"
)

type BrokerServer struct {
	broker broker.Broker
	logger logger.Logger
}

func NewBrokerServer(broker broker.Broker, logger logger.Logger) *BrokerServer {
	return &BrokerServer{
		broker: broker,
		logger: logger,
	}
}

func (s *BrokerServer) Start(ctx context.Context) error {
	opts := s.GetStartOptions()
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	go func() {
		defer func(ctx context.Context) {
			if err := s.broker.Disconnect(); err != nil {
				s.logger.Errorf(ctx, "Failed to shutdown broker server: %v", err)
			}
			s.logger.Info(ctx, "Stop Broker Server")
		}(ctx)

		<-ctx.Done()
	}()

	return nil
}

func HandleBrokerEvent[TReq any, TRes any](b broker.Broker, e broker.Event, opts ...BrokerEventHandlerOption) error {
	var options BrokerEventHandlerOptions
	for _, opt := range opts {
		opt(&options)
	}

	res, err := handleEvent[TReq, TRes](e)

	// Publish to reply topic if present
	replyTopic := options.replyTopic
	if len(replyTopic) != 0 {
		var brokerRes interface{}
		if err != nil {
			brokerRes = broker.FailureResponse(err)
		} else {
			brokerRes = broker.SuccessResponse[TRes](res)
		}

		body, err := json.Marshal(brokerRes)
		if err != nil {
			return err
		}

		fMsg := broker.Message{
			Body: body,
		}

		if e.Message() != nil {
			fMsg.Headers = e.Message().Headers
		}

		return b.Publish(replyTopic, &fMsg)
	}

	return nil
}

func handleEvent[TReq any, TRes any](e broker.Event) (res TRes, err error) {
	// Step 1: Parse the request
	if e.Message() == nil || len(e.Message().Body) == 0 {
		return *new(TRes), broker.EmptyMessageError{}
	}

	body := e.Message().Body
	var request TReq
	err = json.Unmarshal(body, &request)
	if err != nil {
		return *new(TRes), broker.InvalidDataFormatError{}
	}

	// Step 2: Send the request to request pipeline
	res, err = pipeline.Send[TReq, TRes](context.TODO(), request)
	if err != nil {
		return *new(TRes), err
	}

	return res, nil
}

type BrokerEventHandlerOption func(*BrokerEventHandlerOptions)

type BrokerEventHandlerOptions struct {
	replyTopic string
}

func WithReplyTopic(replyTopic string) BrokerEventHandlerOption {
	return func(opts *BrokerEventHandlerOptions) {
		opts.replyTopic = replyTopic
	}
}
