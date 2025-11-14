package schedule

import (
	"fmt"
	"go-event/internal/notification"
	"go-event/internal/participant"
	"go-event/internal/user/repositories"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	repo            Repository
	notifService    notification.Service
	participantRepo participant.Repository
	userRepo        repositories.UserRepository
	cron            *gocron.Scheduler
}

func NewScheduler(
	repo Repository,
	notifService notification.Service,
	participantRepo participant.Repository,
	userRepo repositories.UserRepository,
) *Scheduler {
	return &Scheduler{
		repo:            repo,
		notifService:    notifService,
		participantRepo: participantRepo,
		userRepo:        userRepo,
		cron:            gocron.NewScheduler(time.UTC),
	}
}

func (s *Scheduler) Start() {
	// Jalankan setiap 1 menit untuk cek pending jobs
	s.cron.Every(1).Minute().Do(s.processPendingJobs)
	
	log.Println("Scheduler started - checking jobs every 1 minute")
	s.cron.StartAsync()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) processPendingJobs() {
	jobs, err := s.repo.FindPending()
	if err != nil {
		log.Printf("scheduler: failed to fetch pending jobs: %v", err)
		return
	}

	now := time.Now()
	for _, job := range jobs {
		// Cek apakah waktu run_at sudah lewat
		if job.RunAt.Before(now) || job.RunAt.Equal(now) {
			log.Printf("scheduler: processing job ID %d, type: %s, event: %d", job.ID, job.JobType, job.EventID)
			
			if err := s.executeJob(&job); err != nil {
				log.Printf("scheduler: failed to execute job ID %d: %v", job.ID, err)
				// Update status jadi failed
				if updateErr := s.repo.UpdateStatus(job.ID, StatusFailed); updateErr != nil {
					log.Printf("scheduler: failed to update job status to failed: %v", updateErr)
				}
			} else {
				// Update status jadi done
				if updateErr := s.repo.UpdateStatus(job.ID, StatusDone); updateErr != nil {
					log.Printf("scheduler: failed to update job status to done: %v", updateErr)
				}
				log.Printf("scheduler: job ID %d executed successfully", job.ID)
			}
		}
	}
}

func (s *Scheduler) executeJob(job *ScheduleJob) error {
	switch job.JobType {
	case JobTypeReminder:
		return s.sendReminderNotification(job)
	case JobTypeEndEvent:
		return s.sendEndEventNotification(job)
	default:
		return fmt.Errorf("unknown job type: %s", job.JobType)
	}
}

func (s *Scheduler) sendReminderNotification(job *ScheduleJob) error {
	// Ambil semua participant dari event
	participants, err := s.participantRepo.FindByEventID(job.EventID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) == 0 {
		log.Printf("scheduler: no participants found for event %d", job.EventID)
		return nil
	}

	// Format tanggal event
	eventDate := job.Event.StartTime.Format("02 Jan 2006 15:04")

	// Kirim notifikasi ke setiap participant
	successCount := 0
	for _, p := range participants {
		// Ambil data user untuk email
		userInfo, err := s.userRepo.GetByID(p.UserID)
		if err != nil {
			log.Printf("scheduler: failed to get user %d: %v", p.UserID, err)
			continue
		}

		message := fmt.Sprintf("Reminder: Event '%s' akan dimulai segera pada %s", 
			job.Event.Title, 
			eventDate)

		req := &notification.CreateNotificationRequest{
			UserID:  p.UserID,
			EventID: &job.EventID,
			Type:    notification.NotifReminder,
			Message: message,
		}

		// Kirim notifikasi dengan email
		if _, err := s.notifService.CreateNotificationWithEmail(req, userInfo.Email, userInfo.Name); err != nil {
			log.Printf("scheduler: failed to send reminder to user %d: %v", p.UserID, err)
		} else {
			successCount++
		}
	}

	log.Printf("scheduler: sent %d reminder notifications for event %d", successCount, job.EventID)
	return nil
}

func (s *Scheduler) sendEndEventNotification(job *ScheduleJob) error {
	// Ambil semua participant dari event
	participants, err := s.participantRepo.FindByEventID(job.EventID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) == 0 {
		log.Printf("scheduler: no participants found for event %d", job.EventID)
		return nil
	}

	// Kirim notifikasi ke setiap participant
	successCount := 0
	for _, p := range participants {
		// Ambil data user untuk email
		userInfo, err := s.userRepo.GetByID(p.UserID)
		if err != nil {
			log.Printf("scheduler: failed to get user %d: %v", p.UserID, err)
			continue
		}

		message := fmt.Sprintf("Event '%s' telah selesai. Terima kasih atas partisipasi Anda!", 
			job.Event.Title)

		req := &notification.CreateNotificationRequest{
			UserID:  p.UserID,
			EventID: &job.EventID,
			Type:    notification.NotifUpdate,
			Message: message,
		}

		// Kirim notifikasi dengan email
		if _, err := s.notifService.CreateNotificationWithEmail(req, userInfo.Email, userInfo.Name); err != nil {
			log.Printf("scheduler: failed to send end event notification to user %d: %v", p.UserID, err)
		} else {
			successCount++
		}
	}

	log.Printf("scheduler: sent %d end event notifications for event %d", successCount, job.EventID)
	return nil
}
