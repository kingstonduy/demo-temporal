package app

import (
	"context"
	"errors"
	"fmt"
	shared "kingstonduy/demo-temporal/saga"
	"time"

	"go.temporal.io/sdk/activity"
)

var timeout = time.Second * 1

func ValidateAccount(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Validate Account activity starts")

	url := fmt.Sprintf("http://%s/otp/verify/", shared.NAPAS_SERVICE_URL)
	var responseType shared.NapasEntity
	err := shared.PostApi(url, &input, &responseType)
	if err != nil {
		return errors.New("Cannot validate account")
	}
	log.Info("💡Validate Account activity ends")
	return nil
	// fmt.Println("💡Validate account")
	// time.Sleep(timeout)
	// return nil
}

func UpdateStateCreated(ctx context.Context, input shared.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	err = shared.CreateEntity(db, input)
	if err != nil {
		return errors.New("Cannot create transaction")
	}
	log.Info("💡Update state created successfully")
	return nil

	// fmt.Println("💡Update State created")
	// time.Sleep(timeout)
	// return nil
}

func UpdateStateCreateCompensate(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Rollback update state created starts")
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	err = shared.DeleteEntity(db, input)
	if err != nil {
		return errors.New("Cannot delete transaction")
	}
	log.Info("💡Rollback update state created successfully")
	return nil

	// fmt.Println("💡Rollback update state created")
	// time.Sleep(timeout)
	// return nil
}

func LimitCut(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut activity starts")

	url := fmt.Sprintf(`http://%s/api/v1/account/limit`, shared.LIMITATION_SERVICE_URL)
	var responseType shared.NapasEntity
	err := shared.PostApi(url, &input, &responseType)
	if err != nil {
		return errors.New("Cannot validate account")
	}
	log.Info("💡Validate Account activity ends")
	return nil
	// fmt.Println("💡Limit cut successfully")
	// time.Sleep(timeout)
	// return nil
}

func LimitCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut compensate activity starts")

	url := fmt.Sprintf(`http://%s/api/v1/account/limit`, shared.LIMITATION_SERVICE_URL)
	var responseType shared.NapasEntity
	err := shared.PostApi(url, &input, &responseType)
	if err != nil {
		return errors.New("Cannot validate account")
	}
	log.Info("💡Validate Account activity ends")
	return nil
	// fmt.Println("💡Rollback limit cut")
	// time.Sleep(timeout)
	// return nil
}

func UpdateStateLimitCut(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Update State limit cut successfully")
	time.Sleep(timeout)
	return nil
}

func UpdateStateLimitCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback update state limit cut")
	time.Sleep(timeout)
	return nil
}

func MoneyCut(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Money cut successfully")
	time.Sleep(timeout)
	return nil
}

func MoneyCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback money cut")
	time.Sleep(timeout)
	return nil
}

func UpdateStateMoneyCut(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Update State money cut successfully")
	time.Sleep(timeout)
	return nil
}

func UpdateStateMoneyCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback update state money cut")
	time.Sleep(timeout)
	return nil
}

func UpdateMoney(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Update money successfully")
	time.Sleep(timeout)
	return nil
}

func UpdateMoneyCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback update money")
	time.Sleep(timeout)
	return nil
}

func UpdateStateTransactionCompleted(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Update  state transaction completed")
	return errors.New("")
	// time.Sleep(timeout)
	// return nil
}

func UpdateStateTransactionCompletedCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback update state transaction completed")
	time.Sleep(timeout)
	return nil
}
