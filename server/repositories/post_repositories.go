package repositories

import (
	"waysgallery/models"

	"gorm.io/gorm"
)

type PostRepositories interface {
	FindPosts() ([]models.Post, error)
	GetPostByID(ID int) (models.Post, error)
	CreatePost(post models.Post) (models.Post, error)
	CreatePhoto(photo models.Photo) (models.Photo, error)
	GetPhotoByPostID(ID int) ([]models.Photo, error)
	GetUserPostByID(ID int) (models.User, error)
}

func RepositoryPost(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("Photos").Preload("User").Find(&posts).Error

	return posts, err
}

func (r *repository) GetPostByID(ID int) (models.Post, error) {
	var post models.Post
	err := r.db.Preload("Photos").Preload("User").First(&post, ID).Error

	return post, err
}

func (r *repository) CreatePost(post models.Post) (models.Post, error) {
	err := r.db.Create(&post).Error

	return post, err
}

func (r *repository) CreatePhoto(photo models.Photo) (models.Photo, error) {
	err := r.db.Create(&photo).Error
	return photo, err
}

func (r *repository) GetPhotoByPostID(ID int) ([]models.Photo, error) {
	var photo []models.Photo
	err := r.db.Find(&photo, "post_id", ID).Error
	return photo, err
}

func (r *repository) GetUserPostByID(ID int) (models.User, error) {
	var user models.User
	err := r.db.Preload("Post").Preload("Followings").Preload("Followings.Post").Preload("Followers").Preload("Followers.Post").First(&user, ID).Error

	return user, err
}
