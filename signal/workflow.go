package signal

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type MySignal struct {
	Message string
}

func SignalWorkflow(ctx workflow.Context) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
		HeartbeatTimeout:    time.Minute * 10,
		// RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	signalChan := workflow.GetSignalChannel(ctx, "MySignal")
	var signal string = ""

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		signalChan.Receive(ctx, &signal)
		fmt.Println("💡Received signal", signal)
		_ = workflow.ExecuteActivity(ctx, InputActivity, signal).Get(ctx, nil)
	})

	for i := 0; i < 5; i++ {
		workflow.ExecuteActivity(ctx, BlockingActivity, i+1).Get(ctx, nil)
	}

	// wait all go routines don
	wg.Wait(ctx)
	return nil

}
