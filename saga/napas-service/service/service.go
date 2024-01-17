package service

import (
	"errors"
	shared "kingstonduy/demo-temporal/saga"
	"kingstonduy/demo-temporal/saga/napas-service/repository"
)

func VerifyAccount(input shared.ValidateAccountInput) (shared.NapasEntity, error) {
	db, err := repository.GetConnection()
	if err != nil {
		return shared.NapasEntity{}, errors.New("Cannot connect to database")
	}

	napasEntity, err := repository.GetUserByID(db, input.AccountId)
	if err != nil {
		return shared.NapasEntity{}, errors.New("Cannot find account")
	}

	return napasEntity, nil
}
