package bankdto

type CreateBankDTO struct {
	Name string `json:"name" form:"name" validate:"required"`
}
