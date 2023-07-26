package projectdto

type CreateProjectDTO struct {
	Description string `json:"description" form:"description" validate:"required"`
	OrderID     string `json:"order_id" form:"description" validate:"required"`
}

type UpdateProjectDTO struct {
	Description string `json:"description" form:"description"`
}
