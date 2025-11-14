package routes

import (
	"go-event/internal/user/controllers"
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, ctrl *controllers.UserController, cfg *config.Config) {
	users := app.Group("/api/users")

	// Protected routes - require authentication
	users.Get("/profile", middlewares.Authenticate(cfg), ctrl.GetProfile)
	users.Put("/profile", middlewares.Authenticate(cfg), ctrl.UpdateProfile)
	users.Post("/change-password", middlewares.Authenticate(cfg), ctrl.ChangePassword)

	// Admin only routes
	users.Get("/", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetAllUsers)
	users.Get("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetUserByID)
	users.Delete("/:id", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.DeleteUser)
	users.Get("/role/:role", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.GetUsersByRole)
}
