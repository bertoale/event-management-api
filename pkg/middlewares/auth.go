package middlewares

import (
	"strings"

	"go-event/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah struct untuk JWT payload
type Claims struct {
	ID   uint   `json:"id"`   // User ID dari database
	Role string `json:"role"` // Role as string to avoid import cycle
	jwt.RegisteredClaims
}

// UserLocals is a simplified user struct for context
type UserLocals struct {
	ID    uint
	Name  string
	Email string
	Role  string
}

// ✅ Authenticate Middleware
// Memastikan token JWT valid dan user ada di database
func Authenticate(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari Authorization header atau cookie
		token := c.Get("Authorization")
		if token != "" && strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			token = c.Cookies("token")
		}

		// Jika token tidak ditemukan
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Akses ditolak. Token tidak ditemukan.",
			})
		}

		// Parse dan verifikasi token JWT
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !tkn.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid atau kadaluarsa.",
			})
		}
		// Ambil user dari database
		// We need to query the user to verify they exist
		// But we don't import user package to avoid cycle
		// Instead we check if user exists in DB directly
		db := config.GetDB()
		
		// Query to check if user exists
		var count int64
		if err := db.Table("users").Where("id = ?", claims.ID).Count(&count).Error; err != nil || count == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User tidak ditemukan.",
			})
		}

		// Store user ID and role in context
		c.Locals("userID", claims.ID)
		c.Locals("userRole", claims.Role)

		// Lanjut ke middleware berikutnya
		return c.Next()
	}
}

// ✅ Authorize Middleware
// Mengecek apakah user memiliki salah satu role yang diizinkan
func Authorize(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil user role dari context (harus lewat Authenticate dulu)
		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User belum terautentikasi.",
			})
		}

		// Jika roles kosong, berarti semua user boleh
		if len(roles) == 0 {
			return c.Next()
		}

		// Cek apakah user.role ada di daftar roles yang diizinkan
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		// Jika tidak memiliki izin
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Akses ditolak. Anda tidak memiliki izin yang sesuai.",
		})
	}
}
