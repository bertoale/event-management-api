package participant

import "time"

// EventRepository interface untuk menghindari circular dependency
type EventRepository interface {
GetByID(id uint) (*EventInfo, error)
}

// EventInfo untuk menghindari import event package
type EventInfo struct {
ID          uint
Title       string
Description string
Location    string
StartTime   time.Time
EndTime     time.Time
OrganizerID uint
}
