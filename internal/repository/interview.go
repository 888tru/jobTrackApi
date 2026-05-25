package repository

import (
	"github.com/888tru/jobtrack-api/internal/model"
	"gorm.io/gorm"
)

type InterviewRepository interface {
	Create(interview *model.Interview) error
	FindByApplicationID(appID uint) ([]model.Interview, error)
	FindByID(id, appID uint) (*model.Interview, error)
	Update(interview *model.Interview) error
	Delete(id, appID uint) error
	CountByUserID(userID uint) (int64, error)
}

type interviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db: db}
}

func (r *interviewRepository) Create(interview *model.Interview) error {
	return r.db.Create(interview).Error
}

func (r *interviewRepository) FindByApplicationID(appID uint) ([]model.Interview, error) {
	var interviews []model.Interview
	err := r.db.Where("job_application_id = ?", appID).Order("scheduled_at ASC").Find(&interviews).Error
	return interviews, err
}

func (r *interviewRepository) FindByID(id, appID uint) (*model.Interview, error) {
	var interview model.Interview
	err := r.db.Where("id = ? AND job_application_id = ?", id, appID).First(&interview).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}

func (r *interviewRepository) Update(interview *model.Interview) error {
	return r.db.Save(interview).Error
}

func (r *interviewRepository) Delete(id, appID uint) error {
	return r.db.Where("id = ? AND job_application_id = ?", id, appID).Delete(&model.Interview{}).Error
}

func (r *interviewRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Interview{}).
		Joins("JOIN job_applications ON job_applications.id = interviews.job_application_id").
		Where("job_applications.user_id = ?", userID).
		Count(&count).Error
	return count, err
}
