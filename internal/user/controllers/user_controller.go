package controllers

import (
	"go-event/internal/user"
	"go-event/internal/user/services"
	"go-event/pkg/config"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService services.UserService
	cfg         *config.Config
}

func NewUserController(userService services.UserService, cfg *config.Config) *UserController {
	return &UserController{
		userService: userService,
		cfg:         cfg,
	}
}

// GetProfile - Get current user profile
func (ctrl *UserController) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)

	userResponse, err := ctrl.userService.GetProfile(user.ID)
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
func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := ctrl.userService.GetAllUsers()
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
func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	userResponse, err := ctrl.userService.GetUserByID(uint(userID))
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
func (ctrl *UserController) UpdateProfile(c *fiber.Ctx) error {
	users := c.Locals("user").(*user.User)
	var req user.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	updatedUser, err := ctrl.userService.UpdateProfile(users.ID, &req)
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
func (ctrl *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	if err := ctrl.userService.DeleteUser(uint(userID)); err != nil {
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
func (ctrl *UserController) GetUsersByRole(c *fiber.Ctx) error {
	role := c.Params("role")

	users, err := ctrl.userService.GetUsersByRole(role)
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
func (ctrl *UserController) ChangePassword(c *fiber.Ctx) error {
	users := c.Locals("user").(*user.User)
	var req user.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := ctrl.userService.ChangePassword(users.ID, &req); err != nil {
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