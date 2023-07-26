package models

type Bank struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type BankResponse struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (BankResponse) TableName() string {
	return "banks"
}
