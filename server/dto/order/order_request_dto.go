package orderdto

type CreateOrderDTO struct {
	Title       string `json:"title" form:"title" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	StartDate   string `json:"start_date" form:"start_date" validate:"required"`
	EndDate     string `json:"end_date" form:"end_date" validate:"required"`
	Price       string `json:"price" form:"price" validate:"required"`
	OrderToID   string `json:"order_to_id" form:"order_to_id" validate:"required"`
}

type UpdateOrderDTO struct {
	Status string `json:"status" form:"id" validate:"required"`
}
