package participant

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Register(participant *Participant) error
	FindByEventAndUser(eventID uint, userID uint) (*Participant, error)
	FindByEventID(eventID uint) ([]Participant, error)
	Delete(participant *Participant) error
}

type repository struct {
	db *gorm.DB
}

// Delete implements Repository.
func (r *repository) Delete(participant *Participant) error {
	return r.db.Delete(participant).Error
}

// FindByEventAndUser implements Repository.
func (r *repository) FindByEventAndUser(eventID uint, userID uint) (*Participant, error) {
	var participant Participant
	err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &participant, err
}

// FindByEventID implements Repository.
func (r *repository) FindByEventID(eventID uint) ([]Participant, error) {
	var participants []Participant
	err := r.db.Preload("User").Where("event_id = ?", eventID).Find(&participants).Error
	return participants, err
}

// Register implements Repository.
func (r *repository) Register(participant *Participant) error {
	return r.db.Create(participant).Error
}

func Newrepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
