package postdto

import "waysgallery/models"

type PostResponseDTO struct {
	ID          int                 `json:"id" gorm:"primaryKey"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Photos      []models.Photo      `json:"photos"`
	User        models.UserResponse `json:"user" gorm:"foreignKey:UserID"`
}
