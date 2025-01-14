package gwt

import (
	"github.com/gofiber/fiber/v2"
)

type TokenRequest struct {
	UserLogin  string   `json:"userlogin"`
	UserAccess []string `json:"useraccess"`
}

// Временный обработчик для получения токена
func TokenHandler(c *fiber.Ctx) error {
	var request TokenRequest
	if err := c.BodyParser(&request); err != nil && request.UserLogin != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid or empty request body"})
	}
	// Генерируем токен
	token, err := GenerateToken(request.UserLogin, request.UserAccess)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not generate token"})
	}

	// Возвращаем токен
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Token": token})
}
