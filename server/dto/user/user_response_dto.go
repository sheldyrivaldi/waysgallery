package userdto

import "waysgallery/models"

type UserResponseDTO struct {
	ID         int            `json:"id" gorm:"primaryKey"`
	Fullname   string         `json:"fullname"`
	Email      string         `json:"email" gorm:"unique"`
	Avatar     string         `json:"avatar"`
	Greeting   string         `json:"greeting"`
	Balance    int            `json:"balance"`
	Banner     string         `json:"banner"`
	Followings []*models.User `json:"followings"`
	Followers  []*models.User `json:"followers"`
	Post       []models.Post  `json:"post"`
	Role       string         `json:"role"`
}
