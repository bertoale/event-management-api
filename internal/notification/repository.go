package notification

import "gorm.io/gorm"

type Repository interface {
	Create(notification *Notification) error
	GetByUserID(userID uint) ([]Notification, error)
	MarkAsRead(notificationID uint) error
	Delete(notification *Notification) error
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(notification *Notification) error {
	return r.db.Create(notification).Error
}

// Delete implements Repository.
func (r *repository) Delete(notification *Notification) error {
  return r.db.Delete(notification).Error
}


// GetByUserID implements Repository.
func (r *repository) GetByUserID(userID uint) ([]Notification, error) {
	var notifications []Notification
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("sent_at desc").Find(&notifications).Error
	return notifications, err
}

// MarkAsRead implements Repository.
func (r *repository) MarkAsRead(notificationID uint) error {
	return r.db.Model(&Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func Newrepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
