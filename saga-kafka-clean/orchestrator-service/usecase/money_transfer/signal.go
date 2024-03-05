package usecase

import (
	"orchestrator-service/domain"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type AwaitSignals struct {
	FirstSignalTime      time.Time
	NapasAccountSignal   bool
	NapasAccountResponse domain.SaferResponse
	LimitSignal          bool
	LimitResponse        domain.SaferResponse
	T24Signal            bool
	T24Response          domain.SaferResponse
	NapasMoneySignal     bool
	NapasMoneyResponse   domain.SaferResponse
}

func NewAwaitSignals() domain.AwaitSignals {
	return &AwaitSignals{}
}

// SignalToSignalTimeout is them maximum time between signals
var SignalToSignalTimeout = 30 * time.Second

// FromFirstSignalTimeout is the maximum time to receive all signals
var FromFirstSignalTimeout = 60 * time.Second

// Listen to signals Signal1, Signal2, and Signal3
func (a *AwaitSignals) Listen(ctx workflow.Context) {
	log := workflow.GetLogger(ctx)
	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(workflow.GetSignalChannel(ctx, "napas-account"), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &a.NapasAccountResponse)
			a.NapasAccountSignal = true
			log.Info("Signal Received")
		})

		selector.AddReceive(workflow.GetSignalChannel(ctx, "limit"), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &a.LimitResponse)
			a.LimitSignal = true
			log.Info("Signal Received")
		})

		selector.AddReceive(workflow.GetSignalChannel(ctx, "t24"), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &a.T24Response)
			a.T24Signal = true
			log.Info("Signal Received")
		})

		selector.AddReceive(workflow.GetSignalChannel(ctx, "napas-money"), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &a.NapasMoneyResponse)
			a.NapasMoneySignal = true
			log.Info("Signal Received")
		})

		selector.Select(ctx)
		if a.FirstSignalTime.IsZero() {
			a.FirstSignalTime = workflow.Now(ctx)
		}
	}
}

// GetNextTimeout returns the maximum time allowed to wait for the next signal.
func (a *AwaitSignals) GetNextTimeout(ctx workflow.Context) (time.Duration, error) {
	if a.FirstSignalTime.IsZero() {
		panic("FirstSignalTime is not yet set")
	}
	total := workflow.Now(ctx).Sub(a.FirstSignalTime)
	totalLeft := FromFirstSignalTimeout - total
	if totalLeft <= 0 {
		return 0, temporal.NewApplicationError("FromFirstSignalTimeout", "timeout")
	}
	if SignalToSignalTimeout < totalLeft {
		return SignalToSignalTimeout, nil
	}
	return totalLeft, nil
}
