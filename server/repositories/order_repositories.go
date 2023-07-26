package repositories

import (
	"waysgallery/models"

	"gorm.io/gorm"
)

type OrderRepositories interface {
	FindOrders() ([]models.Order, error)
	FindOrdersByClientID(ID int) ([]models.Order, error)
	FindOrdersByVendorID(ID int) ([]models.Order, error)
	CreateOrder(order models.Order) (models.Order, error)
	GetOrderByID(ID int) (models.Order, error)
	UpdateOrder(order models.Order) (models.Order, error)
	UpdateBalance(user models.User) (string, error)
	GetUserOrderByID(ID int) (models.User, error)
}

func RepositoryOrder(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderBy").Preload("OrderTo").Find(&orders).Error

	return orders, err
}

func (r *repository) FindOrdersByClientID(ID int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderBy").Preload("OrderTo").Find(&orders, "order_by_id = ?", ID).Error

	return orders, err
}

func (r *repository) FindOrdersByVendorID(ID int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderBy").Preload("OrderTo").Find(&orders, "order_to_id = ?", ID).Error

	return orders, err
}

func (r *repository) CreateOrder(order models.Order) (models.Order, error) {
	err := r.db.Create(&order).Error
	return order, err
}

func (r *repository) GetOrderByID(ID int) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderBy").Preload("OrderTo").First(&order, ID).Error

	return order, err
}

func (r *repository) UpdateOrder(order models.Order) (models.Order, error) {
	err := r.db.Save(&order).Error

	return order, err
}

func (r *repository) UpdateBalance(user models.User) (string, error) {
	err := r.db.Save(&user).Error

	return "Success update balance", err
}

func (r *repository) GetUserOrderByID(ID int) (models.User, error) {
	var user models.User
	err := r.db.Preload("Post").First(&user, ID).Error

	return user, err
}
