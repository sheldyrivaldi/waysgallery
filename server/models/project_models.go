package models

type Project struct {
	ID          int            `json:"id"`
	Description string         `json:"description"`
	Photos      []PhotoProject `json:"photos"`
	OrderID     int            `json:"order_id"`
	Order       Order          `json:"-" gorm:"foreignKey:OrderID"`
}

type ProjectResponse struct {
	ID          int            `json:"id"`
	Description string         `json:"description"`
	Photos      []PhotoProject `json:"photos_project"`
	OrderID     int            `json:"order_id"`
	Order       Order          `json:"-"`
}

func (ProjectResponse) TableName() string {
	return "projects"
}
