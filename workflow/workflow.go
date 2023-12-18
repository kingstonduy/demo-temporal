package workflow

import (
	"log"
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
		StartToCloseTimeout: time.Minute * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	futureA := workflow.ExecuteActivity(ctx, activity.RegisterAccount, nil)

	accountB := model.Account{
		Cif: input.Cif1,
	}

	futureB := workflow.ExecuteActivity(ctx, activity.RegisterSms, accountB)
	// futureC := workflow.ExecuteActivity(ctx, activity.RegisterAccount, nil)

	var resultA model.Account
	errA := futureA.Get(ctx, &resultA)
	if errA != nil {
		log.Fatal("Register account failed, err=", errA)
		return errA
	}

	var resultB model.Account
	errB := futureB.Get(ctx, &resultB)
	if errB != nil {
		log.Fatal("Register SMS failed, err=", errA)
		return errA
	}

	return nil
}

func AsyncWorkFlow(ctx workflow.Context) error {
	// register  account
	// register sms
	// send notification to sms
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var account model.Account
	err := workflow.ExecuteActivity(ctx, activity.RegisterAccount, nil).Get(ctx, &account)
	if err != nil {
		log.Fatal("Register account failed, err=", err)
		return err
	}
	log.Println("Register account completely, Account=", account)

	err = workflow.ExecuteActivity(ctx, activity.RegisterSms, &account).Get(ctx, &account)
	if err != nil {
		log.Fatal("Register sms failed, err=", err)
		return err
	}
	log.Println("Register sms completely, Account=", account)

	err = workflow.ExecuteActivity(ctx, activity.NotificationSms, &account).Get(ctx, &account)
	if err != nil {
		log.Fatal("Failed to send notification, err=", err)
		return err
	}
	log.Println("Send notification successfully", account)

	return nil
}
