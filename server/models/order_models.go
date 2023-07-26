package models

type Order struct {
	ID          int          `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	Price       int          `json:"price"`
	OrderByID   int          `json:"order_by_id"`
	OrderBy     UserResponse `json:"order_by" gorm:"foreignKey:OrderByID"`
	OrderToID   int          `json:"order_to_id"`
	OrderTo     UserResponse `json:"order_to" gorm:"foreignKey:OrderToID"`
	Status      string       `json:"status" gorm:"default:'Waiting Accept'"`
}

type OrderResponse struct {
	ID          int          `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	Price       int          `json:"price"`
	OrderByID   int          `json:"order_by_id"`
	OrderBy     UserResponse `json:"order_by" gorm:"-"`
	OrderToID   int          `json:"order_to_id"`
	OrderTo     UserResponse `json:"order_to" gorm:"-"`
	Status      string       `json:"status"`
}

func (Order) TableName() string {
	return "orders"
}

func (OrderResponse) TableName() string {
	return "orders"
}
