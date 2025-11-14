package event

import (
	"go-event/internal/participant"
)

// EventRepositoryAdapter mengadaptasi event.Repository ke participant.EventRepository
type EventRepositoryAdapter struct {
repo Repository
}

// NewEventRepositoryAdapter membuat adapter baru
func NewEventRepositoryAdapter(repo Repository) participant.EventRepository {
return &EventRepositoryAdapter{repo: repo}
}

// GetByID implements participant.EventRepository
func (a *EventRepositoryAdapter) GetByID(id uint) (*participant.EventInfo, error) {
event, err := a.repo.GetByID(id)
if err != nil {
return nil, err
}

return &participant.EventInfo{
ID:          event.ID,
Title:       event.Title,
Description: event.Description,
Location:    event.Location,
StartTime:   event.StartTime,
EndTime:     event.EndTime,
OrganizerID: event.OrganizerID,
}, nil
}
