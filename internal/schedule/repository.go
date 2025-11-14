package schedule

import "gorm.io/gorm"

type Repository interface {
	Create(job *ScheduleJob) error
	FindByEventID(eventID uint) ([]ScheduleJob, error)
	UpdateStatus(id uint, status StatusType) error
	Update(job *ScheduleJob) error
	Delete(id uint) error
	FindPending() ([]ScheduleJob, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) FindPending() ([]ScheduleJob, error) {
	var jobs []ScheduleJob
	err := r.db.
		Where("status = ?", StatusPending).
		Preload("Event").
		Preload("Event.Organizer").
		Find(&jobs).Error

	return jobs, err
}

func (r *repository) Update(job *ScheduleJob) error {
	return r.db.Save(job).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(job *ScheduleJob) error {
	return r.db.Create(job).Error
}

func (r *repository) FindByEventID(eventID uint) ([]ScheduleJob, error) {
	var jobs []ScheduleJob
	err := r.db.Where("event_id = ?", eventID).Find(&jobs).Error
	return jobs, err
}

func (r *repository) UpdateStatus(id uint, status StatusType) error {
	return r.db.Model(&ScheduleJob{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&ScheduleJob{}, id).Error
}
