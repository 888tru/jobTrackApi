package repository

import (
	"github.com/888tru/jobtrack-api/internal/model"
	"gorm.io/gorm"
)

type JobApplicationRepository interface {
	Create(app *model.JobApplication) error
	FindAll(userID uint) ([]model.JobApplication, error)
	FindByID(id, userID uint) (*model.JobApplication, error)
	Update(app *model.JobApplication) error
	Delete(id, userID uint) error
	CountByStatus(userID uint) ([]StatusCount, error)
	CountSince(userID uint, since interface{}) (int64, error)
}

type StatusCount struct {
	Status model.ApplicationStatus
	Count  int64
}

type jobApplicationRepository struct {
	db *gorm.DB
}

func NewJobApplicationRepository(db *gorm.DB) JobApplicationRepository {
	return &jobApplicationRepository{db: db}
}

func (r *jobApplicationRepository) Create(app *model.JobApplication) error {
	return r.db.Create(app).Error
}

func (r *jobApplicationRepository) FindAll(userID uint) ([]model.JobApplication, error) {
	var apps []model.JobApplication
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&apps).Error
	return apps, err
}

func (r *jobApplicationRepository) FindByID(id, userID uint) (*model.JobApplication, error) {
	var app model.JobApplication
	err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("Interviews").First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *jobApplicationRepository) Update(app *model.JobApplication) error {
	return r.db.Save(app).Error
}

func (r *jobApplicationRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.JobApplication{}).Error
}

func (r *jobApplicationRepository) CountByStatus(userID uint) ([]StatusCount, error) {
	var counts []StatusCount
	err := r.db.Model(&model.JobApplication{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&counts).Error
	return counts, err
}

func (r *jobApplicationRepository) CountSince(userID uint, since interface{}) (int64, error) {
	var count int64
	err := r.db.Model(&model.JobApplication{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Count(&count).Error
	return count, err
}
