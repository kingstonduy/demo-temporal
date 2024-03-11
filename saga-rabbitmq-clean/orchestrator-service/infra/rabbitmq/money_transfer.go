package rabbitmq

import (
	"encoding/json"
	"orchestrator-service/bootstrap"
	"orchestrator-service/domain"
	pkg "orchestrator-service/pkg/rabbitmq"
)

type MoneyTransferRabbitMQ struct {
	cfg *bootstrap.Config
}

func NewMoneyTransferRabbitMQ(cfg *bootstrap.Config) domain.MoneyTransferMessageQueue {
	return &MoneyTransferRabbitMQ{
		cfg: cfg,
	}
}

// NapasAccountRequestResponse implements domain.MoneyTransferMessageQueue.
func (mq *MoneyTransferRabbitMQ) NapasAccountRequestResponse(input domain.SaferRequest, queue string) (output domain.NapasAccountResponse, err error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	resStr, err := pkg.RequestAndReply(string(inputBytes), queue, mq.cfg)
	json.Unmarshal([]byte(resStr), &output)
	return output, err
}

// SaferRequestResponse implements domain.MoneyTransferMessageQueue.
func (mq *MoneyTransferRabbitMQ) SaferRequestResponse(input domain.SaferRequest, queue string) (output domain.SaferResponse, err error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	resStr, err := pkg.RequestAndReply(string(inputBytes), queue, mq.cfg)
	json.Unmarshal([]byte(resStr), &output)
	return output, err
}
