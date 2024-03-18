package cron

import (
	"golang.org/x/sync/errgroup"
)

func (s *CronServer) GetStartOptions() []CronServerStartOption {
	return []CronServerStartOption{
		WithCompleteFundTransferBackground(),
	}
}

type CronServerStartOption func(*CronServer) error

func WithCompleteFundTransferBackground() CronServerStartOption {
	return func(cs *CronServer) error {
		eg := new(errgroup.Group)

		eg.Go(func() error {
			// _, err := cs.broker.Subscribe(TopicAccountStatement, func(e broker.Event) error {
			// 	if e.Message() == nil || len(e.Message().Body) == 0 {
			// 		// ignore
			// 		return nil
			// 	}

			// 	body := e.Message().Body
			// 	var request KafkaAccountStatementWrapper
			// 	err := json.Unmarshal(body, &request)
			// 	if err != nil {
			// 		return broker.InvalidDataFormatError{}
			// 	}

			// 	if request.OpType == "I" && request.After != nil {
			// 		pReq := domain.CompleteFundTransferRequest{
			// 			TransNo:    request.After.TransNo,
			// 			TransferAt: time.Now(),
			// 		}
			// 		ctx := context.Background()
			// 		res, err := pipeline.Send[*domain.CompleteFundTransferRequest, *domain.CompleteFundTransferResponse](ctx, &pReq)
			// 		if err != nil {
			// 			cs.logger.Errorf(ctx, "Complete fund transfer failed: %v", err.Error())
			// 		} else {
			// 			cs.logger.Errorf(ctx, "Completed fund transfer: %v", res)
			// 		}
			// 	}
			// 	return nil
			// })
			// return err
			return nil
		})

		return eg.Wait()
	}
}
