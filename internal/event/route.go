package event

import (
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupOrganizerEventRoutes(app *fiber.App, ctrl *Controller, cfg *config.Config) {
	EO := app.Group("/api/events")

	EO.Post("/", middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.CreateEvent)
	EO.Get("/",middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.GetAllEventByUserID)
	EO.Put(":id",middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.UpdateEvent)
	EO.Delete(":id",middlewares.Authenticate(cfg), middlewares.Authorize("organizer"), ctrl.DeleteEvent)
}