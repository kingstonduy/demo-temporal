package workflow

import (
	"encoding/json"
	"fmt"
	"log"
	shared "saga-kafka-notclean/config"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// SignalToSignalTimeout is them maximum time between signals
var SignalToSignalTimeout = 30 * time.Second

// FromFirstSignalTimeout is the maximum time to receive all signals
var FromFirstSignalTimeout = 60 * time.Second

type AwaitSignals struct {
	FirstSignalTime      time.Time
	NapasAccountSignal   bool
	NapasAccountResponse shared.SaferResponse
	LimitSignal          bool
	LimitResponse        shared.SaferResponse
	T24Signal            bool
	T24Response          shared.SaferResponse
	NapasMoneySignal     bool
	NapasMoneyResponse   shared.SaferResponse
}

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

func MoneyTransferWorkflow(ctx workflow.Context, input shared.WorkflowInput) (output shared.WorkflowOutput, err error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 0,
			InitialInterval: time.Second * 5},
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	var compensations Compensations

	defer func() {
		if err != nil {
			// activity failed, and workflow context is canceled
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			compensations.Compensate(disconnectedCtx, true)
		}
	}()

	fmt.Printf("💡Workflow input %+v\n", input)

	saferRequest := shared.SaferRequest{
		WorkflowID:    workflow.GetInfo(ctx).WorkflowExecution.ID,
		RunID:         workflow.GetInfo(ctx).WorkflowExecution.RunID,
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
	}

	var a AwaitSignals
	workflow.Go(ctx, a.Listen)

	err = workflow.ExecuteActivity(ctx, ValidateAccount, saferRequest).Get(ctx, nil)
	if err != nil {
		return
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.NapasAccountSignal
	})
	// Cancellation
	if err != nil || a.NapasAccountResponse.Code != 200 {
		return
	}

	var napasAccountResponse shared.NapasAccountResponse
	err = json.Unmarshal([]byte(a.NapasAccountResponse.Message), &napasAccountResponse)
	if err != nil {
		return
	}

	var transactionEntity = shared.TransactionEntity{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountResponse.AccountName,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Amount:        input.Amount,
		Timestamp:     time.Now().String(),
		State:         "",
	}

	output = shared.WorkflowOutput{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountResponse.AccountName,
		Amount:        input.Amount,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
	}

	err = workflow.ExecuteActivity(ctx, UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(CompensateTransaction, transactionEntity)

	// -------------------------------------------------------------------------------------------------
	err = workflow.ExecuteActivity(ctx, LimitCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.LimitSignal
	})
	// Cancellation
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(LimitCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// -------------------------------------------------------------------------------------------------
	err = workflow.ExecuteActivity(ctx, MoneyCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.T24Signal
	})
	// Cancellation
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(MoneyCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, UpdateStateMoneyCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// -------------------------------------------------------------------------------------------------
	err = workflow.ExecuteActivity(ctx, UpdateMoney, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.NapasMoneySignal
	})
	// Cancellation
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(UpdateMoneyCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// -------------------------------------------------------------------------------------------------
	log.Println("💡Workflow output", output)
	return output, err
}
