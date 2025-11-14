package event

import (
	"errors"
	"fmt"
	"go-event/internal/participant"
	"go-event/internal/user/repositories"
	"go-event/pkg/config"
	"log"

	"gorm.io/gorm"
)

type Service interface {
	CreateEvent(userID uint, req *CreateEventRequest) (*EventResponse, error)
	GetEventByUserID(userID uint) ([]EventResponse, error)
	GetEventByID(eventId uint) (*EventResponse, error)
	UpdateEvent(userID,eventID uint, req *UpdateEventRequest) (*EventResponse, error)
	DeleteEvent(userID,eventID uint) error
}

type service struct {
	repo            Repository
	participantRepo participant.Repository
	userRepo        repositories.UserRepository
	notifService    NotificationService
	cfg             *config.Config
}

// CreateEvent implements Service.
func (s *service) CreateEvent(userID uint, req *CreateEventRequest) (*EventResponse, error) {
	if req.Title == "" || req.Description == "" || req.Location == "" || req.StartTime.IsZero() || req.EndTime.IsZero() {
		return nil, errors.New("all fields are required" )
	}
	event := &Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		OrganizerID: userID,

	}
	
	if err := s.repo.Create(event); err != nil {
		return nil, errors.New("failed to create event: " + err.Error())
	}

	response := &EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		OrganizerID: 	 event.OrganizerID,
	}
	return response, nil

}

// DeleteEvent implements Service.
func (s *service) DeleteEvent(userID uint,eventID uint) error {
	event, err := s.repo.GetByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return errors.New("event not found")
		}
		return errors.New("failed to get event")
	}
	
	if event.OrganizerID != userID {
		return errors.New("unauthorized to delete this event")
	}
	
	// Kirim notifikasi pembatalan ke semua participant (async)
	go func() {
		participants, err := s.participantRepo.FindByEventID(eventID)
		if err != nil {
			log.Printf("Failed to get participants for event %d: %v", eventID, err)
			return
		}
				for _, p := range participants {
			userInfo, err := s.userRepo.GetByID(p.UserID)
			if err != nil {
				log.Printf("Failed to get user %d: %v", p.UserID, err)
				continue
			}
					message := fmt.Sprintf("Event '%s' telah dibatalkan oleh organizer.", event.Title)
			
			if err := s.notifService.SendNotificationWithEmailByString(p.UserID, eventID, "cancellation", message, userInfo.Email, userInfo.Name); err != nil {
				log.Printf("Failed to send cancellation notification to user %d: %v", p.UserID, err)
			}
		}
	}()
	
	if err := s.repo.Delete(event); err != nil {
		return errors.New("failed to delete event")
	}
	return nil
}

// GetEventByID implements Service.
func (s *service) GetEventByID(eventId uint) (*EventResponse, error) {
	task, err := s.repo.GetByID(eventId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, errors.New("failed to get event")
	}
	response := &EventResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Location:    task.Location,
		StartTime:   task.StartTime,
		EndTime:     task.EndTime,
		OrganizerID: task.OrganizerID,
	}
	return response, nil
}

// GetEventByuserID implements Service.
func (s *service) GetEventByUserID(userID uint) ([]EventResponse, error) {
	tasks, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get events")
	}
	var responses []EventResponse
	for _, task := range tasks {
		response := EventResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Location:    task.Location,
			StartTime:   task.StartTime,
			EndTime:     task.EndTime,
			OrganizerID: task.OrganizerID,
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// UpdateEvent implements Service.
func (s *service) UpdateEvent(userID uint,eventID uint, req *UpdateEventRequest) (*EventResponse, error) {
	task, err := s.repo.GetByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, errors.New("failed to get event")
	}
	if task.OrganizerID != userID {
		return nil, errors.New("unauthorized to update this event")
	}

	// Track perubahan untuk notifikasi
	var changes []string
	
	if req.Title != nil && *req.Title != task.Title {
		changes = append(changes, fmt.Sprintf("Judul diubah menjadi: %s", *req.Title))
		task.Title = *req.Title
	}
	if req.Description != nil && *req.Description != task.Description {
		changes = append(changes, "Deskripsi event telah diperbarui")
		task.Description = *req.Description
	}
	if req.Location != nil && *req.Location != task.Location {
		changes = append(changes, fmt.Sprintf("Lokasi diubah menjadi: %s", *req.Location))
		task.Location = *req.Location
	}
	if req.StartTime != nil && !req.StartTime.Equal(task.StartTime) {
		changes = append(changes, fmt.Sprintf("Waktu mulai diubah menjadi: %s", req.StartTime.Format("02 Jan 2006 15:04")))
		task.StartTime = *req.StartTime
	}
	if req.EndTime != nil && !req.EndTime.Equal(task.EndTime) {
		changes = append(changes, fmt.Sprintf("Waktu selesai diubah menjadi: %s", req.EndTime.Format("02 Jan 2006 15:04")))
		task.EndTime = *req.EndTime
	}

	if err := s.repo.Update(task); err != nil {
		return nil, errors.New("failed to update event")
	}
	
	// Kirim notifikasi update ke semua participant jika ada perubahan (async)
	if len(changes) > 0 {
		go func() {
			participants, err := s.participantRepo.FindByEventID(eventID)
			if err != nil {
				log.Printf("Failed to get participants for event %d: %v", eventID, err)
				return
			}
			
			updateMessage := "Perubahan yang dilakukan:\n"
			for _, change := range changes {
				updateMessage += "- " + change + "\n"
			}
					for _, p := range participants {
				userInfo, err := s.userRepo.GetByID(p.UserID)
				if err != nil {
					log.Printf("Failed to get user %d: %v", p.UserID, err)
					continue
				}
						if err := s.notifService.SendNotificationWithEmailByString(p.UserID, eventID, "update", updateMessage, userInfo.Email, userInfo.Name); err != nil {
					log.Printf("Failed to send update notification to user %d: %v", p.UserID, err)
				}
			}
		}()
	}
	
	response := &EventResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Location:    task.Location,
		StartTime:   task.StartTime,
		EndTime:     task.EndTime,
		OrganizerID: task.OrganizerID,
	}
	return response, nil
}

func NewService(repo Repository, participantRepo participant.Repository, userRepo repositories.UserRepository, notifService NotificationService, cfg *config.Config) Service {
	return &service{
		repo:            repo,
		participantRepo: participantRepo,
		userRepo:        userRepo,
		notifService:    notifService,
		cfg:             cfg,
	}
}
