package schedule

import (
	"errors"
	"go-event/internal/event"
	"go-event/pkg/config"
	"time"
)

type Service interface {
	CreateSchedule(req *CreateScheduleRequest) (*ScheduleResponse, error)
	GetSchedulesByEventID(eventID uint) ([]ScheduleResponse, error)
	DeleteSchedule(scheduleID uint, userID uint) error
}

type service struct {
	repo      Repository
	eventRepo event.Repository
	cfg       *config.Config
}

// CreateSchedule implements Service.
func (s *service) CreateSchedule(req *CreateScheduleRequest) (*ScheduleResponse, error) {
	// Validasi event exists
	events, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	// Validasi waktu run_at tidak boleh sebelum waktu sekarang
	if req.RunAt.Before(time.Now()) {
		return nil, errors.New("run_at must be in the future")
	}

	// Validasi waktu run_at tidak boleh setelah event selesai
	if req.RunAt.After(events.EndTime) {
		return nil, errors.New("run_at cannot be after event end time")
	}

	// Buat schedule job
	job := &ScheduleJob{
		EventID:   req.EventID,
		JobType:   req.JobType,
		RunAt:     req.RunAt,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(job); err != nil {
		return nil, errors.New("failed to create schedule: " + err.Error())
	}

	response := &ScheduleResponse{
		ID:      job.ID,
		EventID: job.EventID,
		JobType: job.JobType,
		RunAt:   job.RunAt,
		Status:  job.Status,
	}

	return response, nil
}

// GetSchedulesByEventID implements Service.
func (s *service) GetSchedulesByEventID(eventID uint) ([]ScheduleResponse, error) {
	// Validasi event exists
	_, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	jobs, err := s.repo.FindByEventID(eventID)
	if err != nil {
		return nil, errors.New("failed to retrieve schedules: " + err.Error())
	}

	var responses []ScheduleResponse
	for _, job := range jobs {
		responses = append(responses, ScheduleResponse{
			ID:      job.ID,
			EventID: job.EventID,
			JobType: job.JobType,
			RunAt:   job.RunAt,
			Status:  job.Status,
		})
	}

	return responses, nil
}

// DeleteSchedule implements Service.
func (s *service) DeleteSchedule(scheduleID uint, userID uint) error {
	// Cari schedule berdasarkan ID
	job, err := s.repo.GetByID(scheduleID)
	if err != nil || job == nil {
		return errors.New("schedule not found")
	}

	// Validasi user adalah organizer dari event tersebut
	events, err := s.eventRepo.GetByID(job.EventID)
	if err != nil {
		return errors.New("event not found")
	}

	if events.OrganizerID != userID {
		return errors.New("unauthorized to delete this schedule")
	}

	// Delete schedule
	if err := s.repo.Delete(scheduleID); err != nil {
		return errors.New("failed to delete schedule: " + err.Error())
	}

	return nil
}

func NewService(repo Repository, eventRepo event.Repository, cfg *config.Config) Service {
	return &service{
		repo:      repo,
		eventRepo: eventRepo,
		cfg:       cfg,
	}
}
