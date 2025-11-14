// filepath: d:\CODING\Goevent\internal\notification\service.go
package notification

import (
	"errors"
	"go-event/internal/event"
	"go-event/internal/notification/email"
	"go-event/internal/user"
	"go-event/pkg/config"
	"time"
)

type Service interface {
	CreateNotification(req *CreateNotificationRequest) (*NotificationResponse, error)
	CreateNotificationWithEmail(req *CreateNotificationRequest, userEmail, userName string) (*NotificationResponse, error)
	SendNotificationWithEmail(userID uint, eventID uint, notifType NotifType, message, userEmail, userName string) error
	SendNotificationWithEmailByString(userID uint, eventID uint, notifTypeStr string, message, userEmail, userName string) error
	GetNotificationsByUserID(userID uint) ([]NotificationResponse, error)
	MarkNotificationAsRead(notificationID uint, userID uint) error
	DeleteNotification(notificationID uint, userID uint) error
}

type service struct {
	repo         Repository
	eventRepo    event.Repository
	cfg          *config.Config
	emailService email.Service
}

// CreateNotification implements Service.
func (s *service) CreateNotification(req *CreateNotificationRequest) (*NotificationResponse, error) {
	// Validasi request
	if req.Message == "" {
		return nil, errors.New("message is required")
	}

	// Validasi tipe notifikasi
	validTypes := []NotifType{NotifReminder, NotifUpdate, NotifCancellation}
	isValid := false
	for _, t := range validTypes {
		if req.Type == t {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("invalid notification type")
	}

	notification := &Notification{
		UserID:  req.UserID,
		EventID: req.EventID,
		Type:    req.Type,
		Message: req.Message,
		IsRead:  false,
		SentAt:  time.Now(),
	}

	if err := s.repo.Create(notification); err != nil {
		return nil, errors.New("failed to create notification: " + err.Error())
	}

	response := &NotificationResponse{
		ID:      notification.ID,
		Type:    notification.Type,
		Message: notification.Message,
		IsRead:  notification.IsRead,
		SentAt:  notification.SentAt,
	}

	return response, nil
}

// CreateNotificationWithEmail implements Service.
func (s *service) CreateNotificationWithEmail(req *CreateNotificationRequest, userEmail, userName string) (*NotificationResponse, error) {
	// Buat notifikasi di database terlebih dahulu
	notification, err := s.CreateNotification(req)
	if err != nil {
		return nil, err
	}

	// Kirim email berdasarkan tipe notifikasi (async, tidak block jika gagal)
	go func() {
		var emailErr error
		
		// Ambil detail event jika ada
		var eventTitle string
		var eventDate string
		if req.EventID != nil {
			eventData, err := s.eventRepo.GetByID(*req.EventID)
			if err == nil {
				eventTitle = eventData.Title
				eventDate = eventData.StartTime.Format("02 Jan 2006 15:04")
			} else {
				eventTitle = "Event"
				eventDate = "segera"
			}
		} else {
			eventTitle = "Event"
			eventDate = "segera"
		}

		switch req.Type {
		case NotifReminder:
			emailErr = s.emailService.SendReminderEmail(userEmail, userName, eventTitle, eventDate)
		case NotifCancellation:
			emailErr = s.emailService.SendCancellationEmail(userEmail, userName, eventTitle)
		case NotifUpdate:
			emailErr = s.emailService.SendUpdateEmail(userEmail, userName, eventTitle, req.Message)
		}

		if emailErr != nil {
			// Log error tapi tidak mempengaruhi response
			// Bisa tambahkan logging di sini
		}
	}()

	return notification, nil
}

// SendNotificationWithEmail adalah helper method untuk mengirim notifikasi dari package lain
func (s *service) SendNotificationWithEmail(userID uint, eventID uint, notifType NotifType, message, userEmail, userName string) error {
	req := &CreateNotificationRequest{
		UserID:  userID,
		EventID: &eventID,
		Type:    notifType,
		Message: message,
	}
	
	_, err := s.CreateNotificationWithEmail(req, userEmail, userName)
	return err
}

// SendNotificationWithEmailByString adalah wrapper yang menerima string type untuk digunakan dari package lain
func (s *service) SendNotificationWithEmailByString(userID uint, eventID uint, notifTypeStr string, message, userEmail, userName string) error {
	var notifType NotifType
	
	switch notifTypeStr {
	case "cancellation":
		notifType = NotifCancellation
	case "update":
		notifType = NotifUpdate
	case "reminder":
		notifType = NotifReminder
	default:
		notifType = NotifUpdate
	}
	
	return s.SendNotificationWithEmail(userID, eventID, notifType, message, userEmail, userName)
}

// GetNotificationsByUserID implements Service.
func (s *service) GetNotificationsByUserID(userID uint) ([]NotificationResponse, error) {
	notifications, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve notifications: " + err.Error())
	}

	var responses []NotificationResponse
	for _, notif := range notifications {
		responses = append(responses, NotificationResponse{
			ID:      notif.ID,
			Type:    notif.Type,
			Message: notif.Message,
			IsRead:  notif.IsRead,
			SentAt:  notif.SentAt,
			User: user.UserResponse{
				ID:    notif.User.ID,
				Name:  notif.User.Name,
				Email: notif.User.Email,
				Role:  string(notif.User.Role),
			},
		})
	}

	return responses, nil
}

// MarkNotificationAsRead implements Service.
func (s *service) MarkNotificationAsRead(notificationID uint, userID uint) error {
	// Get notification untuk validasi ownership
	notifications, err := s.repo.GetByUserID(userID)
	if err != nil {
		return errors.New("failed to retrieve notifications")
	}

	// Cek apakah notifikasi milik user
	var found bool
	for _, notif := range notifications {
		if notif.ID == notificationID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("notification not found or unauthorized")
	}

	if err := s.repo.MarkAsRead(notificationID); err != nil {
		return errors.New("failed to mark notification as read: " + err.Error())
	}

	return nil
}

// DeleteNotification implements Service.
func (s *service) DeleteNotification(notificationID uint, userID uint) error {
	// Get notification untuk validasi ownership
	notifications, err := s.repo.GetByUserID(userID)
	if err != nil {
		return errors.New("failed to retrieve notifications")
	}

	// Cari notifikasi yang akan dihapus
	var targetNotif *Notification
	for i := range notifications {
		if notifications[i].ID == notificationID {
			targetNotif = &notifications[i]
			break
		}
	}

	if targetNotif == nil {
		return errors.New("notification not found or unauthorized")
	}

	if err := s.repo.Delete(targetNotif); err != nil {
		return errors.New("failed to delete notification: " + err.Error())
	}
	return nil
}

func NewService(repo Repository, eventRepo event.Repository, emailService email.Service, cfg *config.Config) Service {
	return &service{
		repo:         repo,
		eventRepo:    eventRepo,
		emailService: emailService,
		cfg:          cfg,
	}
}
