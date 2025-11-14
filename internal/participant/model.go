package participant

import (
	"go-event/internal/user"
	"time"
)

type StatusType string

const (
	StatusRegistered  StatusType = "registered"
	StatusAttended  	StatusType = "attended"
	StatusCancelled 	StatusType = "cancelled"
)

// 🧱 Entity (database model)
type Participant struct {
	ID        uint      	`json:"id" gorm:"primaryKey"`
	EventID   uint      	`json:"event_id"`
	UserID    uint      	`json:"user_id"`
	Status    StatusType	`json:"status"`
	CreatedAt time.Time 	`json:"created_at"`

	// Event event.Event `json:"event" gorm:"foreignKey:EventID"` // Removed to avoid circular import
	User  user.User   `json:"user" gorm:"foreignKey:UserID"`
}

// 📩 Request structs
type RegisterParticipantRequest struct {
	EventID uint `json:"event_id" validate:"required"`
	UserID  uint `json:"user_id" validate:"required"`
}

// 📤 Response structs
type ParticipantResponse struct {
	ID      uint              `json:"id"`
	Status  string            `json:"status"`
	User    user.UserResponse `json:"user"`
	EventID uint              `json:"event_id"`
}
