package schedule

import (
	"go-event/internal/event"
	"time"
)

type JobType string
type StatusType string

const (
	JobTypeReminder JobType = "reminder"
	JobTypeEndEvent JobType = "end_event"

	StatusPending StatusType = "pending"
	StatusDone    StatusType = "done"
	StatusFailed  StatusType = "failed"
)

// 🧱 Entity untuk database
type ScheduleJob struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	EventID   uint       `json:"event_id"`
	JobType   JobType    `json:"job_type"`
	RunAt     time.Time  `json:"run_at"`
	Status    StatusType `json:"status"`
	CreatedAt time.Time  `json:"created_at"`

	Event event.Event `json:"event" gorm:"foreignKey:EventID"`
}

// 📩 Request struct
type CreateScheduleRequest struct {
	EventID uint      `json:"event_id" validate:"required"`
	JobType JobType   `json:"job_type" validate:"required,oneof=reminder end_event"`
	RunAt   time.Time `json:"run_at" validate:"required"`
}

// 📤 Response struct
type ScheduleResponse struct {
	ID      uint       `json:"id"`
	EventID uint       `json:"event_id"`
	JobType JobType    `json:"job_type"`
	RunAt   time.Time  `json:"run_at"`
	Status  StatusType `json:"status"`
}
