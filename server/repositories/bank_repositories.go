package repositories

import (
	"waysgallery/models"

	"gorm.io/gorm"
)

type BankRepositories interface {
	CreateBank(bank models.Bank) (models.Bank, error)
	FindBanks() ([]models.Bank, error)
	GetBankByID(ID int) (models.Bank, error)
	DeleteBank(bank models.Bank) (models.Bank, error)
}

func RepositoryBank(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindBanks() ([]models.Bank, error) {
	var banks []models.Bank
	err := r.db.Find(&banks).Error

	return banks, err
}

func (r *repository) CreateBank(bank models.Bank) (models.Bank, error) {
	err := r.db.Create(&bank).Error

	return bank, err
}

func (r *repository) GetBankByID(ID int) (models.Bank, error) {
	var bank models.Bank
	err := r.db.First(&bank, ID).Error

	return bank, err
}

func (r *repository) DeleteBank(bank models.Bank) (models.Bank, error) {
	err := r.db.Delete(&bank).Error

	return bank, err
}
