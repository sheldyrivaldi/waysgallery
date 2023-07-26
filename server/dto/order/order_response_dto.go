package orderdto

import "waysgallery/models"

type OrderResponseDTO struct {
	ID          int                 `json:"id" gorm:"primaryKey"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	StartDate   string              `json:"start_date"`
	EndDate     string              `json:"end_date"`
	Price       int                 `json:"price"`
	OrderByID   int                 `json:"-"`
	OrderBy     models.UserResponse `json:"order_by" gorm:"foreignKey:OrderByID"`
	OrderToID   int                 `json:"-"`
	OrderTo     models.UserResponse `json:"order_to" grom:"foreignKey:OrderToID"`
	Status      string              `json:"status"`
}
