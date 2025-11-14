package controllers

import (
	"go-event/internal/user"
	"go-event/internal/user/services"
	"go-event/pkg/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService 
	cfg *config.Config
}

func NewAuthController(authService services.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		cfg: cfg,
	}
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req user.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})

	}

	token, userResponse, err := ctrl.authService.Login(req)
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

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req user.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	userResponse, err := ctrl.authService.Register(req)
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