package repositories

import (
	"waysgallery/models"

	"gorm.io/gorm"
)

type ProjectRepositories interface {
	CreateProject(project models.Project) (models.Project, error)
	GetProjectByOrderID(ID int) (models.Project, error)
	GetProjectByID(ID int) (models.Project, error)
	UpdateProject(project models.Project) (models.Project, error)
	CreatePhotoProject(photo models.PhotoProject) (models.PhotoProject, error)
	GetPhotoProjectByProjectID(ID int) ([]models.PhotoProject, error)
	DeletePhoto(photo models.PhotoProject) (string, error)
}

func RepositoryProject(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateProject(project models.Project) (models.Project, error) {
	err := r.db.Create(&project).Error

	return project, err
}

func (r *repository) GetProjectByOrderID(ID int) (models.Project, error) {
	var project models.Project
	err := r.db.Preload("Photos").First(&project, "order_id", ID).Error

	return project, err
}

func (r *repository) GetProjectByID(ID int) (models.Project, error) {
	var project models.Project
	err := r.db.Preload("Photos").First(&project, ID).Error

	return project, err
}

func (r *repository) UpdateProject(project models.Project) (models.Project, error) {
	err := r.db.Save(&project).Error

	return project, err
}

func (r *repository) CreatePhotoProject(photo models.PhotoProject) (models.PhotoProject, error) {
	err := r.db.Create(&photo).Error
	return photo, err
}

func (r *repository) GetPhotoProjectByProjectID(ID int) ([]models.PhotoProject, error) {
	var photos []models.PhotoProject
	err := r.db.Find(&photos, "project_id", ID).Error
	return photos, err
}

func (r *repository) DeletePhoto(photo models.PhotoProject) (string, error) {
	err := r.db.Delete(&photo).Error

	return "Delete Success", err
}
