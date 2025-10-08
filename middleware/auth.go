package middleware

import (
	"TM4/utils"
	"strings"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya admin yang diizinkan",
			})
		}
		return c.Next()
	}
}
