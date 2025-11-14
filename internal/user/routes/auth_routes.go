package routes

import (
	"go-event/internal/user/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, ctrl *controllers.AuthController) {
	auth := app.Group("/api/auth")

	auth.Post("/register", ctrl.Register)
	auth.Post("/login", ctrl.Login)
}