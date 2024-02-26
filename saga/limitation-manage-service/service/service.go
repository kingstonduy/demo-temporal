package service

import (
	"errors"
	shared "saga"
)

func LimitService(input shared.SaferRequest) error {
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	var limitEntity = shared.AccountLimitEntity{}
	err = shared.GetEntityByID(db, input.AccountId, &limitEntity)
	if err != nil {
		return errors.New("Cannot find account")
	}

	if limitEntity.Amount < input.Amount {
		return errors.New("Not enough line of credit")
	}

	limitEntity.Amount -= input.Amount
	err = shared.UpdateEntity(db, limitEntity)
	if err != nil {
		return errors.New("Cannot update account")
	}
	return nil
}
