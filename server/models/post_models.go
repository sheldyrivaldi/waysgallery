package models

type Post struct {
	ID          int          `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Photos      []Photo      `json:"photos"`
	UserID      int          `json:"user_id"`
	User        UserResponse `json:"user" gorm:"foreignKey:UserID"`
}

type PostResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Photos      []Photo `json:"photos"`
}

func (PostResponse) TableName() string {
	return "posts"
}
