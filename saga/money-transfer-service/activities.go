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

	url := fmt.Sprintf("http://%s/api/v1/account/verify", shared.NAPAS_SERVICE_HOST_PORT)
	var responseType shared.NapasEntity
	err := shared.PostApi(url, &shared.ValidateAccountInput{input.ToAccountId}, &responseType)
	if err != nil {
		log.Error("🔥Validate Account activity failed")
		return err
	}

	log.Info("💡Validate Account activity successfully")
	return nil
}

func UpdateStateCreated(ctx context.Context, input shared.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")
	db, err := shared.GetConnection()
	if err != nil {
		log.Error("🔥Cannot connect to database")
		return err
	}

	err = shared.CreateEntity(db, input)
	if err != nil {
		log.Error("🔥Cannot create transaction")
		return err
	}

	log.Info("💡Update state created successfully")
	return nil
}

func CompensateTransaction(ctx context.Context, input shared.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Compensate Transaction starts")
	db, err := shared.GetConnection()
	if err != nil {
		log.Error("🔥Cannot connect to database")
		return err
	}

	err = shared.UpdateEntity(db, input)
	if err != nil {
		log.Error("🔥Cannot create compensate transaction")
		return err
	}

	log.Info("💡Compensate Transaction successfully")
	return nil
}

func LimitCut(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	url := fmt.Sprintf("http://%s/api/v1/account/limit", shared.LIMITATION_SERVICE_HOST_PORT)
	var responseType shared.NapasEntity

	log.Info("Limit cut Account activity starts")

	err := shared.PostApi(url,
		&shared.SaferRequest{
			TransactionId: input.TransactionId,
			AccountId:     input.FromAccountId,
			Amount:        input.Amount,
		}, &responseType)
	if err != nil {
		log.Error("🔥Limit cut Account activity failed")
		return err
	}

	log.Info("💡Limit cut Account activity successfully")
	return nil
}

func LimitCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("Limit cut compensate activity starts")
	url := fmt.Sprintf("http://%s/api/v1/account/limit", shared.LIMITATION_SERVICE_HOST_PORT)
	var responseType shared.NapasEntity

	err := shared.PostApi(url,
		&shared.SaferRequest{
			TransactionId: input.TransactionId,
			AccountId:     input.FromAccountId,
			Amount:        -input.Amount, // compensate
		}, &responseType)
	if err != nil {
		log.Error("🔥Limit cut compensate activity failed")
		return err
	}

	log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func UpdateStateLimitCut(ctx context.Context, input shared.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	err = shared.UpdateEntity(db, input)
	if err != nil {
		log.Error("🔥Cannot update state")
		return err
	}
	log.Info("💡Update state successfully")
	return nil
}

func MoneyCut(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut Account activity starts")

	url := fmt.Sprintf("http://%s/api/v1/amount/cut", shared.T24_SERVICE_HOST_PORT)
	var responseType shared.SaferResponse
	err := shared.PostApi(url, &shared.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)
	if err != nil {
		log.Error("🔥Money cut Account activity failed")
		return err
	}

	log.Info("💡Money cut Account activity successfully")
	return nil
}

func MoneyCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut compensate activity starts")

	url := fmt.Sprintf("http://%s/api/v1/amount/add", shared.T24_SERVICE_HOST_PORT)
	var responseType shared.SaferResponse
	input.Amount = -input.Amount

	err := shared.PostApi(url, &shared.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)
	if err != nil {
		log.Error("🔥Money cut compensate activity failed")
		return err
	}

	log.Info("💡Money cut compensate activity successfully")
	return nil
}

func UpdateStateMoneyCut(ctx context.Context, input shared.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	err = shared.UpdateEntity(db, input)
	if err != nil {
		log.Error("🔥Cannot update state")
		return err
	}
	log.Info("💡Update state successfully")
	return nil
}

func UpdateMoney(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("update money napas cut compensate activity starts")

	url := fmt.Sprintf("http://%s/api/v1/account/update", shared.NAPAS_SERVICE_HOST_PORT)
	var responseType shared.NapasEntity
	err := shared.PostApi(url, &shared.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.ToAccountId,
		Amount:        input.Amount,
	}, &responseType)
	if err != nil {
		log.Error("🔥Limit cut compensate activity failed")
		return err
	}
	log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func UpdateMoneyCompensate(ctx context.Context, input shared.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("update money napas cut compensate activity starts")

	url := fmt.Sprintf("http://%s/api/v1/account/update", shared.NAPAS_SERVICE_HOST_PORT)
	input.Amount = -input.Amount
	var responseType shared.NapasEntity
	err := shared.PostApi(url, &shared.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)
	if err != nil {
		log.Error("🔥Limit cut compensate activity failed")
		return err
	}
	log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func UpdateStateTransactionCompleted(ctx context.Context, input shared.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")
	db, err := shared.GetConnection()
	if err != nil {
		log.Error("🔥Cannot connect to database")
		return err
	}

	err = shared.UpdateEntity(db, input)
	if err != nil {
		log.Error("🔥Cannot update state")
		return err
	}
	log.Info("💡Update state  successfully")
	return nil
}
