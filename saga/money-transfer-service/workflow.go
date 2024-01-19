package app

import (
	shared "kingstonduy/demo-temporal/saga"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Compensations struct {
	compensations []any
	arguments     [][]any
}

func (s *Compensations) AddCompensation(activity any, parameters ...any) {
	s.compensations = append(s.compensations, activity)
	s.arguments = append(s.arguments, parameters)
}

func (s Compensations) Compensate(ctx workflow.Context, inParallel bool) {
	if !inParallel {
		// Compensate in Last-In-First-Out order, to undo in the reverse order that activies were applied.
		for i := len(s.compensations) - 1; i >= 0; i-- {
			errCompensation := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...).Get(ctx, nil)
			if errCompensation != nil {
				workflow.GetLogger(ctx).Error("Executing compensation failed", "Error", errCompensation)
			}
		}
	} else {
		selector := workflow.NewSelector(ctx)
		for i := 0; i < len(s.compensations); i++ {
			execution := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...)
			selector.AddFuture(execution, func(f workflow.Future) {
				if errCompensation := f.Get(ctx, nil); errCompensation != nil {
					workflow.GetLogger(ctx).Error("Executing compensation failed", "Error", errCompensation)
				}
			})
		}
		for range s.compensations {
			selector.Select(ctx)
		}

	}
}

func MoneyTransferWorkflow(ctx workflow.Context, info shared.TransactionInfo) (err error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 2},
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

	// just read the database, dont need to compensate
	// compensations.AddCompensation(ValidateAccount)
	err = workflow.ExecuteActivity(ctx, ValidateAccount, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(UpdateStateCreateCompensate)
	err = workflow.ExecuteActivity(ctx, UpdateStateCreated, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(LimitCutCompensate)
	err = workflow.ExecuteActivity(ctx, LimitCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(UpdateStateLimitCutCompensate)
	err = workflow.ExecuteActivity(ctx, UpdateStateLimitCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(MoneyCutCompensate)
	err = workflow.ExecuteActivity(ctx, MoneyCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(UpdateStateMoneyCutCompensate)
	err = workflow.ExecuteActivity(ctx, UpdateStateMoneyCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(UpdateMoneyCompensate)
	err = workflow.ExecuteActivity(ctx, UpdateMoney, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	compensations.AddCompensation(UpdateStateTransactionCompletedCompensate)
	err = workflow.ExecuteActivity(ctx, UpdateStateTransactionCompleted, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
