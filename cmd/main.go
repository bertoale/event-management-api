package main

import (
	"fmt"
	"go-event/internal/event"
	"go-event/internal/notification"
	"go-event/internal/notification/email"
	"go-event/internal/participant"
	"go-event/internal/schedule"
	"go-event/internal/user"

	"go-event/pkg/config"
	"go-event/pkg/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CorsOrigin,
		AllowCredentials: true,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	if err := config.Connect(cfg); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	// Manual migration for vertical architecture
	db := config.GetDB()
	tables := []interface{}{
		&user.User{},
		&event.Event{}, // tambahkan model Event ke migrasi
		&participant.Participant{}, // tambahkan model Event ke migrasi
		&schedule.ScheduleJob{}, // tambahkan model Event ke migrasi
		&notification.Notification{}, // tambahkan model Notification ke migrasi
	}
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("✅ Migrasi database berhasil.")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the REST API",
			"version": "1.0.0",
			"timestamp": fiber.Map{},
		})
	})	// Initialize email service (dibutuhkan untuk semua service yang kirim email)
	emailService := email.NewService(cfg)
	
	// Initialize repositories
	userRepo := user.NewRepository(db)
	eventRepo := event.Newrepository(db)
	participantRepo := participant.Newrepository(db)
	scheduleRepo := schedule.NewRepository(db)
	notificationRepo := notification.Newrepository(db)
	
	// Create adapter for event repository to avoid circular dependency
	eventRepoAdapter := event.NewEventRepositoryAdapter(eventRepo)
	
	// Initialize notification service (dibutuhkan oleh event service dan scheduler)
	notificationService := notification.NewService(notificationRepo, eventRepo, emailService, cfg)
	notificationController := notification.NewController(notificationService, cfg)
	
	// Initialize services
	userService := user.NewService(userRepo, emailService, cfg)
	userController := user.NewController(userService, cfg)
	
	// Initialize event service (dengan dependency notification untuk update/cancel)
	eventService := event.NewService(eventRepo, participantRepo, userRepo, notificationService, cfg)
	eventController := event.NewController(eventService, cfg)
	
	// Initialize participant service (dengan email service untuk konfirmasi registrasi)
	participantService := participant.NewService(participantRepo, eventRepoAdapter, userRepo, emailService, cfg)
	participantController := participant.NewController(participantService, *cfg)
	
	// Initialize schedule service
	scheduleService := schedule.NewService(scheduleRepo, eventRepo, cfg)
	scheduleController := schedule.NewController(scheduleService, cfg)
	// Initialize scheduler with all dependencies
	scheduler := schedule.NewScheduler(scheduleRepo, notificationService, participantRepo, userRepo)
	scheduler.Start()
	defer scheduler.Stop()

	// Use vertical layer routes
	user.SetupUserRoutes(app, userController, cfg)
	event.SetupOrganizerEventRoutes(app, eventController, cfg)
	participant.SetupParticipantRoute(app, participantController, cfg)
	schedule.SetupScheduleRoutes(app, scheduleController, cfg)
	notification.SetupNotificationRoutes(app, notificationController, cfg)

	app.Use(middlewares.NotFound)

	port := cfg.Port
	log.Printf("Server is running on port %s", port)
	log.Printf("Local: http://localhost:%s", port)
	log.Printf("Environment: %s", cfg.NodeEnv)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}