// filepath: d:\CODING\Goevent\internal\notification\controller.go
package notification

import (
	"go-event/pkg/config"
	"strconv"

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

	// Bisa parse JSON & x-www-form-urlencoded
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
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
	userID := c.Locals("userID").(uint)

	notifications, err := ctrl.service.GetNotificationsByUserID(userID)
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
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	notificationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid notification ID",
		})
	}

	err = ctrl.service.MarkNotificationAsRead(uint(notificationID), userID)
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
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	notificationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid notification ID",
		})
	}

	err = ctrl.service.DeleteNotification(uint(notificationID), userID)
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
