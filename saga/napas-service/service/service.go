package service

import (
	"errors"
	shared "kingstonduy/demo-temporal/saga"
)

func VerifyAccount(input shared.ValidateAccountInput) (shared.NapasEntity, error) {
	db, err := shared.GetConnection()
	if err != nil {
		return shared.NapasEntity{}, errors.New("Cannot connect to database")
	}

	var napasEntity = shared.NapasEntity{}
	err = shared.GetEntityByID(db, input.AccountId, &napasEntity)
	if err != nil {
		return shared.NapasEntity{}, errors.New("Cannot find account")
	}

	return napasEntity, nil
}

func UpdateMoney(input shared.SaferRequest) error {
	db, err := shared.GetConnection()
	if err != nil {
		return errors.New("Cannot connect to database")
	}

	var napasEntity = shared.NapasEntity{}
	err = shared.GetEntityByID(db, input.AccountId, &napasEntity)
	if err != nil {
		return errors.New("Cannot find account")
	}

	napasEntity.Amount += input.Amount
	if napasEntity.Amount < 0 {
		return errors.New("Not enough money")
	}

	err = shared.UpdateEntity(db, napasEntity)
	if err != nil {
		return errors.New("Cannot update account")
	}

	return nil
}
