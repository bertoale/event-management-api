// filepath: d:\CODING\Goevent\internal\notification\controller.go
package notification

import (
	"go-event/internal/user"
	"go-event/pkg/config"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
	cfg     *config.Config
}

func NewController(service Service, cfg *config.Config) *Controller {
	return &Controller{
		service: service,
		cfg:     cfg,
	}
}

func (ctrl *Controller) CreateNotification(c *fiber.Ctx) error {
	var req CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi tipe notifikasi
	validTypes := []string{string(NotifReminder), string(NotifUpdate), string(NotifCancellation)}
	isValid := false
	for _, t := range validTypes {
		if strings.ToLower(string(req.Type)) == t {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid notification type. must be 'reminder', 'update', or 'cancellation'",
		})
	}

	notification, err := ctrl.service.CreateNotification(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "notification created successfully",
		"notification": notification,
	})
}

func (ctrl *Controller) GetNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)

	notifications, err := ctrl.service.GetNotificationsByUserID(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "notifications retrieved successfully",
		"notifications": notifications,
	})
}

func (ctrl *Controller) MarkAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)
	id := c.Params("id")

	notificationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid notification ID",
		})
	}

	err = ctrl.service.MarkNotificationAsRead(uint(notificationID), user.ID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "notification not found or unauthorized" {
			statusCode = fiber.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "notification marked as read successfully",
	})
}

func (ctrl *Controller) DeleteNotification(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)
	id := c.Params("id")

	notificationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid notification ID",
		})
	}

	err = ctrl.service.DeleteNotification(uint(notificationID), user.ID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "notification not found or unauthorized" {
			statusCode = fiber.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "notification deleted successfully",
	})
}
