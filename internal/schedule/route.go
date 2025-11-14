package schedule

import (
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupScheduleRoutes(app *fiber.App, ctrl *Controller, cfg *config.Config) {
	schedules := app.Group("/api/events")

	schedules.Post("/:id/schedules", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.CreateSchedule)
	schedules.Get("/:id/schedules", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.GetSchedules)
	schedules.Delete("/schedules/:id", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.DeleteSchedule)
}