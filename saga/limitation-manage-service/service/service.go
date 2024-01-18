package service

import (
	"errors"
	shared "kingstonduy/demo-temporal/saga"
)

func LimitService(input shared.SaferRequest) error {
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	var limitEntity = shared.AccountLimitEntity{}
	err = shared.GetUserByID(db, input.AccountId, &limitEntity)
	if err != nil {
		return errors.New("Cannot find account")
	}

	if limitEntity.Amount < input.Amount {
		return errors.New("Not enough line of credit")
	}

	limitEntity.Amount -= input.Amount
	err = shared.UpdateUser(db, limitEntity)
	if err != nil {
		return errors.New("Cannot update account")
	}
	return nil
}
