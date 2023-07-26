package models

type PhotoProject struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	ProjectID int     `json:"project_id"`
	Project   Project `json:"-"`
	URL       string  `json:"url"`
}

type PhotoProjectResponse struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	ProjectID int    `json:"project_id"`
	URL       string `json:"url"`
}

func (PhotoProjectResponse) TableName() string {
	return "photo_projects"
}
