package repository

import (
	"context"
	"errors"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

// Create implements domain.TransactionRepository.
func (repo *transactionRepository) Create(c context.Context, transaction domain.TransactionEntity) error {
	result := repo.db.Create(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete implements domain.TransactionRepository.
func (repo *transactionRepository) DeleteById(c context.Context, transaction domain.TransactionEntity) error {
	result := repo.db.Delete(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID implements domain.TransactionRepository.
func (repo *transactionRepository) GetByID(c context.Context, id string) (domain.TransactionEntity, error) {
	var result domain.TransactionEntity
	err := repo.db.Where("account_id = ?", id).First(&result).Error
	if err != nil {
		return result, errors.New("Cannot find account")
	}
	return result, nil
}

// Update implements domain.TransactionRepository.
func (repo *transactionRepository) Update(c context.Context, transaction domain.TransactionEntity) error {
	result := repo.db.Save(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
