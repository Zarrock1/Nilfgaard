package privelegies

import (
	"context"
	"core_mod/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func checkSelfAtemptS(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int

	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users u JOIN atempts a ON u.id = a.user_id WHERE u.login = $1 AND a.id = $2", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func checkSelfAtemptSP(c *fiber.Ctx) bool {
	t_id, err := strconv.Atoi(c.Params("t_id"))
	if err != nil {
		return false
	}
	u_id, err := strconv.Atoi(c.Params("u_id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users u JOIN disciplines d ON u.id = d.prepod_id JOIN tests t ON t.discipline_id = d.id WHERE u.login = $1 AND t.id = $2", userlogin, t_id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users  WHERE login = $1 AND id = $2", userlogin, u_id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func CreateAtempt(c *fiber.Ctx) error {
	return c.Next()
}
func UpdateAtempt(c *fiber.Ctx) error {
	if checkSelfAtemptS(c) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})

}
func CompleteAtempt(c *fiber.Ctx) error {
	if checkSelfAtemptS(c) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})

}
func GetAtempts(c *fiber.Ctx) error {
	const right = "test:answer:read"
	if checkSelfAtemptSP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})

}
