package service

import (
	shared "kingstonduy/demo-temporal/saga"
	"kingstonduy/demo-temporal/saga/napas-service/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitService(t *testing.T) {
	db, err := repository.GetConnection()
	if err != nil {
		t.Error(err)
	}

	entity := shared.AccountLimitEntity{AccountId: "1", Amount: 1000}
	err = repository.CreateEntity(db, entity)
	if err != nil {
		t.Error(err)
	}
	defer db.Delete(entity)

	input := shared.SaferRequest{AccountId: "1", Amount: 1000}
	err = LimitService(input)
	if err != nil {
		t.Error(err)
	}

	entity1 := shared.AccountLimitEntity{}
	err = repository.GetUserByID(db, "1", &entity1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, entity.Amount, entity1.Amount+input.Amount)
}
