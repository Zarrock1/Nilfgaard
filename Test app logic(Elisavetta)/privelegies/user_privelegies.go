package privelegies

import (
	"context"
	"core_mod/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func checkRights(c *fiber.Ctx, right string) bool {
	arr := c.Locals("access").([]string)

	for _, i := range arr {
		if i == right {
			return true
		}
	}
	return false
}

func checkSelfUser(c *fiber.Ctx) bool {
	idurl, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var idsql int
	err = db.Pool.QueryRow(context.Background(), "SELECT id FROM users WHERE login = $1", userlogin).Scan(&idsql)
	if err != nil {
		return false
	}
	if idurl == idsql {
		return true
	}
	return false
}

func UsersHandler(c *fiber.Ctx) error {
	const right = "user:list:read"

	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func UserHandler(c *fiber.Ctx) error {
	return c.Next()
}

func UserUpdate(c *fiber.Ctx) error {
	const right = "user:fullName:write"
	if checkSelfUser(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func UserTests(c *fiber.Ctx) error {
	const right = "user:data:read"
	if checkSelfUser(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func UserRoles(c *fiber.Ctx) error {
	const right = "user:roles:read"

	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func UserUpdateRoles(c *fiber.Ctx) error {
	const right = "user:roles:write"

	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func UserStatus(c *fiber.Ctx) error {
	const right = "user:block:read"

	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func UserUpdateStatus(c *fiber.Ctx) error {
	const right = "user:block:write"

	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
