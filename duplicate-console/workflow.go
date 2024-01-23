package duplicate_console

import (
	"fmt"
	"log"
	"time"

	model "kingstonduy/demo-temporal/duplicate-console/model"

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

	futureA := workflow.ExecuteActivity(ctx, RegisterAccount, nil)
	futureB := workflow.ExecuteActivity(ctx, RegisterSms, accountB)
	futureC := workflow.ExecuteActivity(ctx, GetBalanceById, input.Cif2)

	var resultA model.Account
	errA := futureA.Get(ctx, &resultA)
	if errA != nil {
		log.Fatal("Register account failed, err=", errA)
		return errA
	}
	fmt.Println("💡Register new account, account=", resultA)

	var resultB model.Account
	errB := futureB.Get(ctx, &resultB)
	if errB != nil {
		log.Fatal("Register SMS failed, err=", errA)
		return errA
	}
	fmt.Println("💡Register sms for account id=", input.Cif1, "account=", resultB)

	var resultC float64
	errC := futureC.Get(ctx, &resultC)
	if errC != nil {
		log.Fatal("Get balance by id failed, err=", errC)
		return errC
	}

	fmt.Println("💡Get balance for account id=", input.Cif2, "balance=", resultC)

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
	err := workflow.ExecuteActivity(ctx, RegisterAccount, nil).Get(ctx, &account)
	if err != nil {
		log.Fatal("Register account failed, err=", err)
		return err
	}
	log.Println("💡Register account successfully, Account=", account)

	_ = workflow.ExecuteActivity(ctx, RegisterSms, &account).Get(ctx, &account)
	log.Println("💡Register sms successfully, Account=", account)

	err = workflow.ExecuteActivity(ctx, NotificationSms, &account).Get(ctx, &account)
	if err != nil {
		log.Fatal("Failed to send notification, err=", err)
		return err
	}
	log.Println("💡🎇The Account id=", account.Cif, "have register SMS notification successfully!")

	return nil
}
