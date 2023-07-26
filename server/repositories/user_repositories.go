package repositories

import (
	"waysgallery/models"

	"gorm.io/gorm"
)

type UserRepositories interface {
	FindUsers() ([]models.User, error)
	GetUserByID(ID int) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	FollowingUser(currentUser, followingUser models.User) (models.User, error)
	FollowedByUser(currentUser, followingUser models.User) (models.User, error)
}

func RepositoryUser(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Post").Preload("Followings").Preload("Followers").Find(&users).Error

	return users, err
}
func (r *repository) GetUserByID(ID int) (models.User, error) {
	var user models.User
	err := r.db.Preload("Post").Preload("Followings").Preload("Followings.Post").Preload("Followers").Preload("Followers.Post").First(&user, ID).Error

	return user, err
}
func (r *repository) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error

	return user, err
}

func (r *repository) FollowingUser(currentUser, followingUser models.User) (models.User, error) {
	err := r.db.Model(&currentUser).Association("Followings").Append(&followingUser)

	return currentUser, err
}

func (r *repository) FollowedByUser(currentUser, followingUser models.User) (models.User, error) {
	err := r.db.Model(&followingUser).Association("Followers").Append(&currentUser)

	return followingUser, err
}
