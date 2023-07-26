package models

type Photo struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	PostID int    `json:"post_id"`
	URL    string `json:"url"`
}

type PhotoResponse struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	PostID int    `json:"post_id"`
	Post   Post   `json:"-"`
	URL    string `json:"url"`
}

func (PhotoResponse) TableName() string {
	return "photos"
}
