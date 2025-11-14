package participant

import (
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupParticipantRoute(app *fiber.App, ctrl *Controller, cfg *config.Config) {
	PR := app.Group("/api/events/:id")

	PR.Post("/join", middlewares.Authenticate(cfg), ctrl.RegisterParticipant)
	PR.Delete("/join",middlewares.Authenticate(cfg), ctrl.CancelParticipant)
	PR.Get(":id",middlewares.Authenticate(cfg),  ctrl.GetParticipant)
}