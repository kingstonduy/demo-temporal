package service

import (
	"errors"
	shared "kingstonduy/demo-temporal/saga"
	"kingstonduy/demo-temporal/saga/napas-service/repository"
)

func LimitService(input shared.SaferRequest) error {
	db, err := repository.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	var limitEntity = shared.AccountLimitEntity{}
	err = repository.GetUserByID(db, input.AccountId, &limitEntity)
	if err != nil {
		return errors.New("Cannot find account")
	}

	if limitEntity.Amount < input.Amount {
		return errors.New("Not enough line of credit")
	}

	limitEntity.Amount -= input.Amount
	err = repository.UpdateUser(db, limitEntity)
	if err != nil {
		return errors.New("Cannot update account")
	}
	return nil
}
