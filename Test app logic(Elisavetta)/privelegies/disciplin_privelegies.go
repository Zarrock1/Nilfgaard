package privelegies

import (
	"context"
	"core_mod/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func checkSelfDisciplinP(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int
	err = db.Pool.QueryRow(context.Background(), "SELECT count(*) FROM users JOIN disciplines ON users.id = disciplines.prepod_id WHERE users.login = $1 AND disciplines.id = $2", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func checkSelfDisciplinS(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int

	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users JOIN users_disciplines ON users.id = users_disciplines.user_id WHERE users.login = $1 AND users_disciplines.discipline_id = $2 and users_disciplines.user_id = users.id", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func checkSelfDisciplinN(c *fiber.Ctx) bool {
	idurl, err := strconv.Atoi(c.Params("s_id"))
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
func DisciplinsHandler(c *fiber.Ctx) error {
	return c.Next()
}
func DisciplinHandler(c *fiber.Ctx) error {
	return c.Next()
}

func DisciplinUpdate(c *fiber.Ctx) error {
	const right = "course:info:write"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinTests(c *fiber.Ctx) error {
	const right = "course:testList"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkSelfDisciplinS(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinTestStatus(c *fiber.Ctx) error {
	const right = "course:test:read"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkSelfDisciplinS(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinTestStatusUpdate(c *fiber.Ctx) error {
	const right = "course:test:write"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinTestCreaite(c *fiber.Ctx) error {
	const right = "course:test:add"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinTestDelete(c *fiber.Ctx) error {
	const right = "course:test:del"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinStudents(c *fiber.Ctx) error {
	const right = "course:userList"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinStudentAdd(c *fiber.Ctx) error {
	const right = "course:user:add"
	if checkSelfDisciplinN(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinStudentDelete(c *fiber.Ctx) error {
	const right = "course:user:del"
	if checkSelfDisciplinN(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}

func DisciplinCreate(c *fiber.Ctx) error {
	const right = "course:add"
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
func DisciplinDeleted(c *fiber.Ctx) error {
	const right = "course:del"
	if checkSelfDisciplinP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
