package middlewares

import (
	"strings"

	"go-event/internal/user"
	"go-event/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah struct untuk JWT payload
type Claims struct {
	ID uint `json:"id"` // User ID dari database
	Role user.RoleType `json:"role"`
	jwt.RegisteredClaims
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
		var user user.User
		if err := config.DB.First(&user, claims.ID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User tidak ditemukan.",
			})
		}

		// Simpan user di context
		c.Locals("user", &user)

		// Lanjut ke middleware berikutnya
		return c.Next()
	}
}

// ✅ Authorize Middleware
// Mengecek apakah user memiliki salah satu role yang diizinkan
func Authorize(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil user dari context (harus lewat Authenticate dulu)
		user, ok := c.Locals("user").(*user.User)
		if !ok || user == nil {
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
			if string(user.Role) == role {
				return c.Next()
			}
		}

		// Jika tidak memiliki izin
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Akses ditolak. Anda tidak memiliki izin yang sesuai.",
		})
	}
}
