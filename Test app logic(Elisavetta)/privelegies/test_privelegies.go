package privelegies

import (
	"context"
	"core_mod/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func checkSelfTestP(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users u JOIN disciplines d ON u.id = d.prepod_id JOIN tests t ON t.discipline_id = d.id WHERE u.login = $1 AND t.id = $2", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func checkSelfTestS(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int

	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users JOIN users_disciplines ud ON users.id = ud.user_id JOIN tests t ON t.discipline_id = ud.discipline_id  WHERE users.login = $1 AND t.id = $2 and ud.user_id = users.id", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func DeletedQuestionFromTest(c *fiber.Ctx) error {
	const right = "test:quest:del"
	if checkSelfTestP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func AddQuestionToTest(c *fiber.Ctx) error {
	const right = "test:quest:add"
	if checkSelfTestP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func ChangeQuestionOrderInTest(c *fiber.Ctx) error {
	const right = "test:quest:update"
	if checkSelfTestP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func GetUsersPassedTest(c *fiber.Ctx) error {
	const right = "test:answer:read"
	if checkSelfTestP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func GetUserMarksTest(c *fiber.Ctx) error {
	const right = "test:answer:read"
	if checkSelfTestP(c) {
		return c.Next()
	}
	if checkSelfTestS(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func GetUserAnswersTest(c *fiber.Ctx) error {
	const right = "test:answer:read"
	if checkSelfTestP(c) {
		return c.Next()
	}
	if checkSelfTestS(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
