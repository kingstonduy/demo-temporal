package repository

import (
	"orchestrator-service/bootstrap"
	"orchestrator-service/domain"

	"go.temporal.io/sdk/temporal"
	"gorm.io/gorm"
)

type MoneytransferRepository struct {
	DB *gorm.DB
}

func NewMoneytransferRepository(cfg *bootstrap.Config) domain.MoneyTransferRepository {
	return &MoneytransferRepository{
		DB: bootstrap.GetDB(cfg),
	}
}

// Create implements domain.MoneyTransferRepository.
func (r *MoneytransferRepository) Create(entity domain.TransactionEntity) error {
	err := r.DB.Create(entity).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}

// Save implements domain.MoneyTransferRepository.
func (r *MoneytransferRepository) Save(entity domain.TransactionEntity) error {
	err := r.DB.Save(entity).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}
