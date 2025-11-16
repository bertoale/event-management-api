package schedule

import (
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupScheduleRoutes(app *fiber.App, ctrl *Controller, cfg *config.Config) {
	schedules := app.Group("/api/schedule/event")
	schedules.Post("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.CreateSchedule)
	schedules.Get("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.GetSchedules)

	schedules2 := app.Group("/api/schedule")
	schedules2.Delete("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.DeleteSchedule)
}