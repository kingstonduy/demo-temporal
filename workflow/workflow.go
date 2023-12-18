package workflow

import (
	"log"
	"time"

	activity "demo-temporal/activity"
	model "demo-temporal/model"

	"go.temporal.io/sdk/workflow"
)

func ParallelWorkFlow(ctx workflow.Context) error {
	// register sms
	// register email
	// withdraw money
	return nil
}

func AsyncWorkFlow(ctx workflow.Context) (model.Account, error) {
	// register  account
	// register sms
	// send notification to sms
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 45,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var account model.Account
	err := workflow.ExecuteActivity(ctx, activity.RegisterAccount, nil).Get(ctx, &account)
	if err != nil {
		log.Fatal("Register account failed, err=", err)
		return model.Account{}, err
	}

	log.Println("Register account completely, Account=", account)

	err = workflow.ExecuteActivity(ctx, activity.RegisterSms, &account).Get(ctx, &account)
	if err != nil {
		log.Fatal("Register sms failed, err=", err)
		return model.Account{}, err
	}

	log.Println("Register sms completely, Account=", account)

	return account, nil
}
