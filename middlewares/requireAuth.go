package middlewares

import (
	"strings"

	"github.com/Investorharry19/voxa-golang-server/utils"
	"github.com/gofiber/fiber/v2"
)

func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, claims, err := utils.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// You can store token claims in Locals (for later handlers)

	c.Locals("userId", claims["id"])

	return c.Next()
}
