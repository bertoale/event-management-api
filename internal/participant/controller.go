package participant

import (
	"go-event/pkg/config"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
	cfg 	 config.Config
}

func NewController(service Service, cfg config.Config) *Controller {
	return &Controller{
		service: service,
		cfg: cfg,
	}
}

func (ctrl *Controller) RegisterParticipant(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	eventID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid event ID",
		})
	}
	
	req := RegisterParticipantRequest{
		EventID: uint(eventID),
		UserID:  userID,
	}

	participant, err := ctrl.service.RegisterParticipant(&req)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "participant already registered" {
			statusCode = fiber.StatusConflict
		} else if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "user not found" {
			statusCode = fiber.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "participant registered successfully",
		"participant": participant,
	})

}

func (ctrl *Controller) CancelParticipant(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")
	eventID, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid event ID",
		})
	}
	err = ctrl.service.CancelParticipant(uint(eventID), userID)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "participant not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "unauthorized to cancel this participant" {
			statusCode = fiber.StatusUnauthorized
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "participant cancelled successfully",
	})
}

func (ctrl *Controller) GetParticipant(c *fiber.Ctx) error {
	userRole := c.Locals("userRole").(string)
	id := c.Params("id")
	eventID, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid event ID",
		})
	}
	
	if userRole != "admin" && userRole != "organizer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized to view participants",
		})
	}

	participants, err := ctrl.service.GetParticipantsByEventID(uint(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "participants retrieved successfully",
		"participants": participants,
	})
}

