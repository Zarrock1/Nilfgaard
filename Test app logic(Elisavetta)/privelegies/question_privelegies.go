package privelegies

import (
	"context"
	"core_mod/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func checkSelfQuestionP(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int
	err = db.Pool.QueryRow(context.Background(), "SELECT count(*) FROM users JOIN questions q ON users.id = q.avtor_id WHERE users.login = $1 AND q.id = $2", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func checkSelfQuestionS(c *fiber.Ctx) bool {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return false
	}

	userlogin := c.Locals("login").(string)
	var count int
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users JOIN atempts a ON users.id = a.user_id JOIN atempts_questions_answers aqa ON a.id =  aqa.atempt_id JOIN questions_versions qv ON qv.id =  aqa.question_version_id WHERE users.login = $1 AND qv.question_id = $2", userlogin, id).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func QuestionsHendler(c *fiber.Ctx) error {
	return c.Next()
	// Нифига не работает и фиг сделаешь в данной логике
}

func QuestionHendler(c *fiber.Ctx) error {
	const right = "quest:read"
	if checkSelfQuestionP(c) {
		return c.Next()
	}
	if checkSelfQuestionS(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}

func QuestionUpdate(c *fiber.Ctx) error {
	const right = "quest:update"
	if checkSelfQuestionP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}

func QuestionCreate(c *fiber.Ctx) error {
	const right = "quest:create"
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}

func QuestionsDelete(c *fiber.Ctx) error {
	const right = "quest:del"
	if checkSelfQuestionP(c) {
		return c.Next()
	}
	if checkRights(c, right) {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"Message": "Нет прав доступа"})
}
