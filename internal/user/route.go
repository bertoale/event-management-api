package user

import (
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, ctrl *Controller, cfg *config.Config) {
	auth := app.Group("/api/auth")
	auth.Post("/register", ctrl.Register)
	auth.Post("/login", ctrl.Login)

	user := app.Group("/api/user")
	user.Get("/profile", middlewares.Authenticate(cfg), ctrl.GetProfile)
	user.Put("/profile", middlewares.Authenticate(cfg), ctrl.UpdateProfile)
	user.Post("/change-password", middlewares.Authenticate(cfg), ctrl.ChangePassword)
	// Admin only routes
	user.Get("/", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetAllUsers)
	user.Get("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetUserByID)
	user.Delete("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.DeleteUser)
	user.Get("/role/:role", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetUsersByRole)
}