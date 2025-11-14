package participant

import (
	"errors"
	"go-event/internal/notification/email"
	"go-event/internal/user"
	"go-event/internal/user/repositories"
	"go-event/pkg/config"
	"log"
	"time"
)

type Service interface {
	RegisterParticipant(req *RegisterParticipantRequest) (*ParticipantResponse, error)
	CancelParticipant(eventID uint, userID uint) error
	GetParticipantsByEventID(eventID uint) ([]ParticipantResponse, error)
}

type service struct {
	repo         Repository
	cfg          *config.Config
	eventRepo    EventRepository
	userRepo     repositories.UserRepository
	emailService email.Service
}

// CancelParticipant implements Service.
func (s *service) CancelParticipant(eventID uint, userID uint) error {
	participant, err := s.repo.FindByEventAndUser(eventID, userID)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.New("Participant not found")
	}
	return s.repo.Delete(participant)
}

// GetParticipantsByEventID implements Service.
func (s *service) GetParticipantsByEventID(eventID uint) ([]ParticipantResponse, error) {
	participants, err := s.repo.FindByEventID(eventID)
	if err != nil {
		return nil, err
	}

	var response []ParticipantResponse
	for _, p := range participants {
		response = append(response, ParticipantResponse{
			ID: p.ID,
			Status: string(p.Status),
			EventID: p.EventID,
			User: user.UserResponse{
				ID:    p.User.ID,
				Name:  p.User.Name,
				Email: p.User.Email,
				Role:  string(p.User.Role),
			},
		})
	}
	return response,  nil
}




// RegisterParticipant implements Service.
func (s *service) RegisterParticipant(req *RegisterParticipantRequest) (*ParticipantResponse, error) {
	events, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}
	existing, err := s.repo.FindByEventAndUser(req.EventID, req.UserID)
	if err != nil {
		return nil, err	}
	
	if existing != nil {
		return nil, errors.New("user already registered for this event")
	}

	participant := &Participant{
		EventID: 		events.ID,
		UserID:  		req.UserID,
		Status: 		StatusRegistered,
		CreatedAt: 	time.Now(),
	}
	
	if err := s.repo.Register(participant); err != nil{
		return nil, errors.New("failed to register participant: " + err.Error())
	}

	users, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, errors.New("failed to get user info")
	}
	
	// Kirim email konfirmasi pendaftaran (async, tidak block jika gagal)
	go func() {
		eventDate := events.StartTime.Format("02 Jan 2006 15:04")
		if err := s.emailService.SendRegistrationConfirmationEmail(
			users.Email, 
			users.Name, 
			events.Title, 
			eventDate, 
			events.Location,
		); err != nil {
			log.Printf("Failed to send registration confirmation email to %s: %v", users.Email, err)
		}
	}()
	
	response := &ParticipantResponse{
		ID:         participant.ID,
		Status:			string(participant.Status),
		User: 			*users.ToResponse(),
		EventID: 		participant.EventID,
	}
	return response, nil
	
}

	


func NewService(repo Repository, eventRepo EventRepository, userRepo repositories.UserRepository, emailService email.Service, cfg *config.Config) Service {
	return &service{
		repo:         repo,
		cfg:          cfg,
		eventRepo:    eventRepo,
		userRepo:     userRepo,
		emailService: emailService,
	}
}
