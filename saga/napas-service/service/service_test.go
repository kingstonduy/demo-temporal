package service

import (
	shared "kingstonduy/demo-temporal/saga"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	db, err := shared.GetConnection()
	if err != nil {
		t.Error(err)
	}

	entity := shared.NapasEntity{AccountId: "1", AccountName: "Duy", Amount: 1000}

	shared.CreateEntity(db, entity)
	if err != nil {
		t.Error(err)
	}

	entity1, err := VerifyAccount(shared.ValidateAccountInput{AccountId: "1"})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, entity.AccountId, entity1.AccountId)
	assert.Equal(t, entity.AccountName, entity1.AccountName)
	assert.Equal(t, entity.Amount, entity1.Amount)

	defer db.Delete(entity)
}

func TestUpdateMoney(t *testing.T) {
	db, err := shared.GetConnection()
	if err != nil {
		t.Error(err)
	}

	entity := shared.NapasEntity{AccountId: "1", AccountName: "Duy", Amount: 1000}

	shared.CreateEntity(db, entity)
	if err != nil {
		t.Error(err)
	}

	err = UpdateMoney(shared.SaferRequest{AccountId: "1", Amount: 2000})
	if err != nil {
		t.Error(err)
	}

	var entity1 shared.NapasEntity
	err = db.Where("account_id = ?", "1").First(&entity1).Error
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, entity.AccountId, entity1.AccountId)
	assert.Equal(t, entity.AccountName, entity1.AccountName)
	assert.NotEqual(t, entity.Amount, entity1.Amount)
	assert.Equal(t, entity.Amount+2000, entity1.Amount)

	defer db.Delete(entity)
}
