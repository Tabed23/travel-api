package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tabed23/travel-api/utils"
)

func AuthMiddleware(c *fiber.Ctx) error {
	//tokenString := c.Get("Authorization")
	tokenString := utils.ExtractToken(c)
	token, err := utils.ParseToken(tokenString)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token")
	}

	c.Set("role", token.Role)

	return c.Next()
}
