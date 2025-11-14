package schedule

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
	return &Controller{service: service, cfg: cfg}
}

func (ctrl *Controller) CreateSchedule(c *fiber.Ctx) error {
	id := c.Params("id")

	eventID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid event ID",
		})
	}

	var req CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Set event ID dari URL params
	req.EventID = uint(eventID)

	// Validasi job type
	validJobTypes := []string{string(JobTypeReminder), string(JobTypeEndEvent)}
	isValid := false
	for _, jt := range validJobTypes {
		if strings.ToLower(string(req.JobType)) == jt {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid job type. must be 'reminder' or 'end_event'",
		})
	}

	schedule, err := ctrl.service.CreateSchedule(&req)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "schedule created successfully",
		"schedule": schedule,
	})
}

func (ctrl *Controller) GetSchedules(c *fiber.Ctx) error {
	id := c.Params("id")

	eventID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid event ID",
		})
	}

	schedules, err := ctrl.service.GetSchedulesByEventID(uint(eventID))
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "schedules retrieved successfully",
		"schedules": schedules,
	})
}

func (ctrl *Controller) DeleteSchedule(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)
	id := c.Params("id")

	scheduleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid schedule ID",
		})
	}

	err = ctrl.service.DeleteSchedule(uint(scheduleID), user.ID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "schedule not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "event not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "unauthorized to delete this schedule" {
			statusCode = fiber.StatusForbidden
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "schedule deleted successfully",
	})
}

