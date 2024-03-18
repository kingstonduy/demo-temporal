package broker

import (
	"github.com/lengocson131002/go-clean-core/transport/broker"
	"github.com/ocb/mcs-money-transfer/domain/account"
	"golang.org/x/sync/errgroup"
)

func (s *BrokerServer) GetStartOptions() []BrokerServerStartOption {
	return []BrokerServerStartOption{
		WithHandlerSubscriptionRoutes(),
	}
}

type BrokerServerStartOption func(*BrokerServer) error

func WithHandlerSubscriptionRoutes() BrokerServerStartOption {
	return func(b *BrokerServer) error {
		eg := new(errgroup.Group)
		eg.Go(func() error {
			_, e := b.broker.Subscribe(TopicRequestCheckBalance, func(e broker.Event) error {
				return HandleBrokerEvent[*account.CheckBalanceRequest, *account.CheckBalanceResponse](b.broker, e, WithReplyTopic(TopicReplyCheckBalance))
			})
			return e
		})

		return eg.Wait()
	}
}
