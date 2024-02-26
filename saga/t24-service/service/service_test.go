package service

import (
	shared "saga"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountService(t *testing.T) {
	db, err := shared.GetConnection()
	if err != nil {
		t.Error(err)
	}

	entity := shared.T24Entity{AccountId: "1", Amount: 1000}
	err = shared.CreateEntity(db, entity)
	if err != nil {
		t.Error(err)
	}
	defer db.Delete(entity)

	input := shared.SaferRequest{AccountId: "1", Amount: 1000}
	err = AmountService(input)
	if err != nil {
		t.Error(err)
	}

	entity1 := shared.T24Entity{}
	err = shared.GetEntityByID(db, "1", &entity1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, entity.Amount+input.Amount, entity1.Amount)
}
