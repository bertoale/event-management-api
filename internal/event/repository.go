package event

import "gorm.io/gorm"

type Repository interface {
	Create(event *Event) error
	GetByID(id uint) (*Event, error)
	Update(event *Event) error
	Delete(event *Event) error
	GetAllByUserID(userID uint ) ([]*Event, error)
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(event *Event) error {
	return r.db.Create(event).Error
}

// Delete implements Repository.
func (r *repository) Delete(event *Event) error {
	return r.db.Delete(event).Error
}

// GetAll implements Repository.
func (r *repository) GetAllByUserID(userID uint) ([]*Event, error) {
	var events []*Event
	if err := r.db.
		Where("organizer_id = ?", userID).
		Order("created_at desc").
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// GetByID implements Repository.
func (r *repository) GetByID(id uint) (*Event, error) {
	var event Event
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// Update implements Repository.
func (r *repository) Update(event *Event) error {
	return r.db.Save(event).Error
}

func Newrepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
