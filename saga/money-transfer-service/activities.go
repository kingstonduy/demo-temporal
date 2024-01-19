package app

import (
	"context"
	"errors"
	"fmt"
	shared "kingstonduy/demo-temporal/saga"
	"time"
)

var timeout = time.Second * 5

func ValidateAccount(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Validate Account successfully")
	time.Sleep(timeout)
	return nil
}

func UpdateStateCreated(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Update State created successfully")
	time.Sleep(timeout)
	return nil
}

func UpdateStateCreateCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback update state created")
	time.Sleep(timeout)
	return nil
}

func LimitCut(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Limit cut successfully")
	time.Sleep(timeout)
	return nil
}

func LimitCutCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback limit cut")
	time.Sleep(timeout)
	return nil
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
	return nil
}

func UpdateStateTransactionCompletedCompensate(ctx context.Context, input shared.TransactionInfo) error {
	fmt.Println("💡Rollback update state transaction completed")
	time.Sleep(timeout)
	return nil
}
