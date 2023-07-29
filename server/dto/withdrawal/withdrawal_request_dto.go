package withdrawaldto

type CreateWithdrawalDTO struct {
	BankID        string `json:"bank_id" form:"bank_id" validate:"required"`
	AccountNumber string `json:"account_number" form:"account_number" validate:"required"`
	Amount        string `json:"amount" form:"amount" validate:"required"`
}

type UpdateWithdrawalDTO struct {
	Status string `json:"status" form:"status" validate:"required"`
}
