package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	db, err := GetConnection()
	if err != nil {
		t.Error(err)
	}

	entity := NapasEntity{AccountId: "1", AccountName: "Duy", Amount: 1000}

	CreateEntity(db, entity)
	if err != nil {
		t.Error(err)
	}

	var entity1 NapasEntity
	err = GetEntityByID(db, entity.AccountId, &entity1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, entity.AccountId, entity1.AccountId)
	assert.Equal(t, entity.AccountName, entity1.AccountName)
	assert.Equal(t, entity.Amount, entity1.Amount)

	defer db.Delete(entity)
}
