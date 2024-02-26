package service

import (
	"errors"
	model "saga-kafka-notclean/money-transfer-service/shared"
)

func LimitService(input model.SaferRequest) error {
	db, err := model.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	var limitEntity = model.AccountLimitEntity{}
	err = model.GetEntityByID(db, input.AccountId, &limitEntity)
	if err != nil {
		return errors.New("Cannot find account")
	}

	if limitEntity.Amount < input.Amount {
		return errors.New("Not enough line of credit")
	}

	limitEntity.Amount -= input.Amount
	err = model.UpdateEntity(db, limitEntity)
	if err != nil {
		return errors.New("Cannot update account")
	}
	return nil
}
