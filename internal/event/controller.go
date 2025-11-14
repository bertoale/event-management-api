package event

import (
	"go-event/pkg/config"
	"strconv"

	"github.com/gofiber/fiber/v2"
)
type Controller struct {
	service Service
	cfg     *config.Config
}

func NewController(service Service, cfg *config.Config ) *Controller {
	return &Controller{
		service: service,
		cfg: cfg,
	}
}

func (ctrl *Controller) CreateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	eventResponse, err := ctrl.service.CreateEvent(userID, &req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "event created successfully",
		"event":   eventResponse,
	})
}

func (ctrl *Controller) GetAllEventByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	events, err := ctrl.service.GetEventByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve events",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "events retrieved successfully",
		"events": events,
	})


}

func (ctrl *Controller) UpdateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	eventId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid event id",
		})
	}

	var req UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	updatedEvent, err := ctrl.service.UpdateEvent(userID, uint(eventId), &req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		}else if err.Error()== "unauthorized to update this event"{
			statusCode = fiber.StatusUnauthorized
		}else if err.Error()== "invalid event data"{
			statusCode = fiber.StatusBadRequest
		}
		return  c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	"message": "event updated successfully",
	"event":   updatedEvent,
})

}

func (ctrl *Controller) DeleteEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	eventId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid event id",
		})
	}
	err = ctrl.service.DeleteEvent(userID, uint(eventId))
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		}else if err.Error()== "unauthorized to delete this event"{
			statusCode = fiber.StatusUnauthorized
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "event deleted successfully",
	})
}