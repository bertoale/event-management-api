// Package middleware contains custom middleware functions
package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler adalah custom error handler untuk Fiber
// Function ini akan dipanggil saat terjadi error di aplikasi
// Function ini di-set di Fiber config saat inisialisasi app
// Parameters:
//   - c: Fiber context
//   - err: Error yang terjadi
// Returns: error (selalu nil karena sudah di-handle)
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code 500 Internal Server Error
	code := fiber.StatusInternalServerError

	// Jika error adalah Fiber error, gunakan status code dari error tersebut
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Return JSON response dengan error message
	// Stack trace tidak ditampilkan untuk security (production)
	return c.Status(code).JSON(fiber.Map{
		"message": err.Error(),
		"stack":   nil, // Don't expose stack trace in production
	})
}

// NotFound adalah handler untuk 404 Not Found
// Handler ini dipanggil untuk semua routes yang tidak terdefinisi
// Function ini di-register sebagai fallback handler di main.go
// Parameters:
//   - c: Fiber context
// Returns: error
func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "NOT FOUND - " + c.OriginalURL(),
	})
}
