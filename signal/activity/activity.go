package activity

import (
	"context"
	"demo-temporal/services"
)

func BlockingActivity(ctx context.Context, input int) error {
	services.PrintService(input)

	return nil
}

func InputActivity(ctx context.Context, input string) error {
	services.InputService(input)

	return nil
}
