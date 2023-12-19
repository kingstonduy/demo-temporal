package workflow

import (
	"time"

	activity "demo-temporal/activity"
	model "demo-temporal/model"

	"go.temporal.io/sdk/workflow"
)

// transfer
func ParallelWorkFlow(ctx workflow.Context, input model.ParallelWorkflowInput) error {
	// create new account
	// register sms existing account
	// check balance existing account
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 15,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	accountB := model.Account{
		Cif: input.Cif1,
	}

	futureA := workflow.ExecuteActivity(ctx, activity.RegisterAccount, nil)
	futureB := workflow.ExecuteActivity(ctx, activity.RegisterSms, accountB)
	futureC := workflow.ExecuteActivity(ctx, activity.GetBalanceById, input.Cif2)

	var resultA model.Account
	errA := futureA.Get(ctx, &resultA)
	if errA != nil {
		return errA
	}

	var resultB model.Account
	errB := futureB.Get(ctx, &resultB)
	if errB != nil {
		return errA
	}

	var resultC float64
	errC := futureC.Get(ctx, &resultC)
	if errC != nil {
		return errC
	}

	return nil
}

func AsyncWorkFlow(ctx workflow.Context) error {
	// register  account
	// register sms
	// send notification to sms
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 15,
		HeartbeatTimeout:    time.Second * 15,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var account model.Account
	err := workflow.ExecuteActivity(ctx, activity.RegisterAccount, nil).Get(ctx, &account)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.RegisterSms, &account).Get(ctx, &account)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.NotificationSms, &account).Get(ctx, &account)
	if err != nil {
		return err
	}

	return nil
}
