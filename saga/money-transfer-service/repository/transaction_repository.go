package repository

import (
	"context"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"

	"gorm.io/gorm"
)

type transactionRepository struct {
	database *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{
		database: db,
	}
}

// Create implements domain.TransactionRepository.
func (*transactionRepository) Create(c context.Context, transaction *domain.TransactionInfo) error {
	panic("unimplemented")
}

// Fetch implements domain.TransactionRepository.
func (*transactionRepository) Fetch(c context.Context) ([]domain.TransactionInfo, error) {
	panic("unimplemented")
}

// GetByEmail implements domain.TransactionRepository.
func (*transactionRepository) GetByEmail(c context.Context, email string) (domain.TransactionInfo, error) {
	panic("unimplemented")
}

// GetByID implements domain.TransactionRepository.
func (*transactionRepository) GetByID(c context.Context, id string) (domain.TransactionInfo, error) {
	panic("unimplemented")
}
