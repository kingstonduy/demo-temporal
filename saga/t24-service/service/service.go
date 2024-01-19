package service

import (
	"errors"
	shared "kingstonduy/demo-temporal/saga"
)

func AmountService(input shared.SaferRequest) error {
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	var t24Entity = shared.T24Entity{}
	err = shared.GetUserByID(db, input.AccountId, &t24Entity)
	if err != nil {
		return errors.New("Cannot find account")
	}

	if t24Entity.Amount+input.Amount < 0 {
		return errors.New("Not enough line of credit")
	}

	t24Entity.Amount += input.Amount
	err = shared.UpdateUser(db, t24Entity)
	if err != nil {
		return errors.New("Cannot update account")
	}
	return nil
}
