package user

import (
	"go-event/pkg/config"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
	cfg         *config.Config
}

func NewController(authService Service, cfg *config.Config) *Controller {
	return &Controller{
		service: authService,
		cfg: cfg,
	}
}


func (ctrl *Controller) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})

	}

	token, userResponse, err := ctrl.service.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   ctrl.cfg.NodeEnv == "production",
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "Login successfully.",
		"token":   token,
		"user":    userResponse,
	})
}

func (ctrl *Controller) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	userResponse, err := ctrl.service.Register(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully.",
		"user":    userResponse,
	})
}




// GetProfile - Get current user profile
func (ctrl *Controller) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	userResponse, err := ctrl.service.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "profile retrieved successfully",
		"user":    userResponse,
	})
}

// GetAllUsers - Get all users (Admin only)
func (ctrl *Controller) GetAllUsers(c *fiber.Ctx) error {
	users, err := ctrl.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "users retrieved successfully",
		"users":   users,
	})
}

// GetUserByID - Get user by ID (Admin only)
func (ctrl *Controller) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	userResponse, err := ctrl.service.GetUserByID(uint(userID))
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user retrieved successfully",
		"user":    userResponse,
	})
}

// UpdateProfile - Update current user profile
func (ctrl *Controller) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	updatedUser, err := ctrl.service.UpdateProfile(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "profile updated successfully",
		"user":    updatedUser,
	})
}

// DeleteUser - Delete user by ID (Admin only)
func (ctrl *Controller) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	if err := ctrl.service.DeleteUser(uint(userID)); err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}

// GetUsersByRole - Get users by role (Admin only)
func (ctrl *Controller) GetUsersByRole(c *fiber.Ctx) error {
	role := c.Params("role")

	users, err := ctrl.service.GetUsersByRole(role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "users retrieved successfully",
		"users":   users,
	})
}

// ChangePassword - Change current user password
func (ctrl *Controller) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := ctrl.service.ChangePassword(userID, &req); err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "invalid old password" {
			statusCode = fiber.StatusUnauthorized
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "password changed successfully",
	})
}