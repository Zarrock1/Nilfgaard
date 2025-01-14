package gwt

import (
	"context"
	"core_mod/db"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Protected — middleware для защиты эндпоинтов с JWT
func Protected(c *fiber.Ctx) error {
	// Извлекаем токен из заголовка Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Токен должен быть в формате "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка подписи с использованием нашего секретного ключа
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Записываем информацию о пользователе из токена в контекст
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	c.Locals("login", claims.UserLogin)   // Добавляем логин пользователя в контекст
	c.Locals("access", claims.UserAccess) // Добовляем права доступа в контекст
	var idsql int
	var count int
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE login = $1", claims.UserLogin).Scan(&count)
	if err != nil {
		log.Println("В базе нет такого user", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token login или ошибка в базе"})
	}
	if count == 1 {
		err = db.Pool.QueryRow(context.Background(), "SELECT id FROM users WHERE login = $1", claims.UserLogin).Scan(&idsql)
		if err != nil {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token login или ошибка в базе"})
		}
	}
	if count == 0 {
		err = db.Pool.QueryRow(context.Background(), "INSERT INTO users (name, blocked, login) VALUES('Аноним',false, $1) RETURNING id", claims.UserLogin).Scan(&idsql)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token login или ошибка в базе"})
		}
	}
	c.Locals("user_id", idsql)
	return c.Next()
}
