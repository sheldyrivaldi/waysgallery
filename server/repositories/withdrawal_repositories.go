package repositories

import (
	"waysgallery/models"

	"gorm.io/gorm"
)

type WithdrawalRepositories interface {
	FindWithdrawals() ([]models.Withdrawal, error)
	GetWithdrawalByID(ID int) (models.Withdrawal, error)
	CreateWithdrawal(withdrawal models.Withdrawal) (models.Withdrawal, error)
	UpdateWithdrawal(withdrawal models.Withdrawal) (models.Withdrawal, error)
	GetUserWithdrawalByID(ID int) (models.User, error)
	UpdateUserWithdrawal(user models.User) (models.User, error)
	FindWithdrawalsByUserID(ID int) ([]models.Withdrawal, error)
}

func RepositoryWithdrawal(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindWithdrawals() ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	err := r.db.Preload("Bank").Preload("User").Find(&withdrawals).Error

	return withdrawals, err
}

func (r *repository) FindWithdrawalsByUserID(ID int) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	err := r.db.Where("user_id = ?", ID).Preload("Bank").Preload("User").Find(&withdrawals).Error

	return withdrawals, err
}

func (r *repository) GetWithdrawalByID(ID int) (models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	err := r.db.Preload("Bank").Preload("User").First(&withdrawal, ID).Error

	return withdrawal, err
}

func (r *repository) CreateWithdrawal(withdrawal models.Withdrawal) (models.Withdrawal, error) {
	err := r.db.Create(&withdrawal).Error

	return withdrawal, err
}

func (r *repository) UpdateWithdrawal(withdrawal models.Withdrawal) (models.Withdrawal, error) {
	err := r.db.Save(&withdrawal).Error

	return withdrawal, err
}

func (r *repository) GetUserWithdrawalByID(ID int) (models.User, error) {
	var user models.User
	err := r.db.Preload("Post").First(&user, ID).Error

	return user, err
}

func (r *repository) UpdateUserWithdrawal(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error

	return user, err
}
