package models

type Withdrawal struct {
	ID            int          `json:"id"`
	UserID        int          `json:"user_id"`
	User          UserResponse `json:"user" gorm:"foreignKey:UserID"`
	BankID        int          `json:"bank_id"`
	Bank          BankResponse `json:"bank"`
	AccountNumber int          `json:"account_number"`
	Amount        int          `json:"amount"`
	Status        string       `json:"status" gorm:"default:'Pending'"`
}

type WithdrawalResponse struct {
	ID            int          `json:"id"`
	UserID        int          `json:"user_id"`
	User          UserResponse `json:"user" gorm:"foreignKey:UserID"`
	BankID        int          `json:"bank_id"`
	Bank          BankResponse `json:"bank" gorm:"foreignKey:BankID"`
	AccountNumber int          `json:"account_number"`
	Amount        int          `json:"amount"`
	Status        string       `json:"status"`
}

func (WithdrawalResponse) TableName() string {
	return "withdrawals"
}
