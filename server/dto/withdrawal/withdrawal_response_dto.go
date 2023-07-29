package withdrawaldto

import "waysgallery/models"

type WithdrawalResponseDTO struct {
	ID            int                 `json:"id"`
	User          models.UserResponse `json:"user"`
	Bank          models.BankResponse `json:"bank"`
	AccountNumber int                 `json:"account_number"`
	Amount        int                 `json:"amount"`
	Status        string              `json:"status"`
}
