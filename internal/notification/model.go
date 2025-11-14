package notification

import (
	"go-event/internal/user"
	"time"
)

type NotifType string

const (
	NotifReminder     NotifType = "reminder"
	NotifUpdate       NotifType = "update"
	NotifCancellation NotifType = "cancellation"
)


// 🧱 Entity (database model)
type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	EventID   *uint     `json:"event_id"`
	Type      NotifType    `json:"type"` // reminder, update, cancellation
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	SentAt    time.Time `json:"sent_at"`
	User      user.User `json:"user" gorm:"foreignKey:UserID"`
}

// 📩 Request structs
type CreateNotificationRequest struct {
	UserID  uint      `json:"user_id" validate:"required"`
	EventID *uint     `json:"event_id"`
	Type    NotifType `json:"type" validate:"required"`
	Message string    `json:"message" validate:"required"`
}


// 📤 Response structs
type NotificationResponse struct {
	ID      uint              `json:"id"`
	Type    NotifType         `json:"type"`
	Message string            `json:"message"`
	IsRead  bool              `json:"is_read"`
	SentAt  time.Time         `json:"sent_at"`
	User    user.UserResponse `json:"user"`
}

