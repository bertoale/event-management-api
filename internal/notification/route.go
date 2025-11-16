// filepath: d:\CODING\Goevent\internal\notification\route.go
package notification

import (
	"go-event/pkg/config"
	"go-event/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes(app *fiber.App, ctrl *Controller, cfg *config.Config) {
	notif := app.Group("/api/notification")

	// Semua user yang authenticated bisa mengakses notifikasi mereka
	notif.Get("/", middlewares.Authenticate(cfg), ctrl.GetNotifications)
	notif.Put("/:id/read", middlewares.Authenticate(cfg), ctrl.MarkAsRead)
	notif.Delete("/:id", middlewares.Authenticate(cfg), ctrl.DeleteNotification)

	// Hanya admin yang bisa create notifikasi (untuk testing atau manual trigger)
	notif.Post("/", middlewares.Authenticate(cfg), middlewares.Authorize("admin"), ctrl.CreateNotification)
}
