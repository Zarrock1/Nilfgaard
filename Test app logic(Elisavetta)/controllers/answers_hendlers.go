package controllers

import (
	"context"
	"core_mod/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UpdateAnswer(c *fiber.Ctx) error {
	var question_version_id int
	var answer_number int
	var id_atempt int
	var answer_id int
	answer_number = -1
	id_question, err := strconv.Atoi(c.Params("q_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}

	student_id := c.Locals("user_id")
	err = db.Pool.QueryRow(context.Background(), "SELECT a.id FROM atempts a JOIN atempts_questions_answers aqa ON a.id = aqa.atempt_id JOIN questions_versions qv ON qv.id = aqa.question_version_id WHERE a.user_id = $1 AND qv.question_id = $2", student_id, id_question).Scan(&id_atempt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с параметрами видимо"})
	}

	if err := c.BodyParser(&answer_number); err != nil || answer_number == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустой "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT aqa.question_version_id, an.id FROM atempts a JOIN tests t ON t.id = a.test_id JOIN atempts_questions_answers aqa ON a.id = aqa.atempt_id JOIN questions_versions qv ON aqa.question_version_id = qv.id JOIN answers an ON an.question_version_id = aqa.question_version_id WHERE t.deleted = false AND a.id = $1 AND qv.question_id = $2 AND an.number = $3", id_atempt, id_question, answer_number).Scan(&question_version_id, &answer_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с параметрами видимо"})
	}
	if question_version_id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с параметрами видимо"})
	}
	_, err = db.Pool.Exec(context.Background(), "UPDATE  atempts_questions_answers SET answer_id = $1 WHERE atempt_id = $2 AND question_version_id = $3", answer_id, id_atempt, question_version_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с чем то..."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}

func DeleteAnswer(c *fiber.Ctx) error {
	var id_atempt int
	var question_version_id int
	id_question, err := strconv.Atoi(c.Params("q_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	student_id := c.Locals("user_id")
	err = db.Pool.QueryRow(context.Background(), "SELECT a.id FROM atempts a JOIN atempts_questions_answers aqa ON a.id = aqa.atempt_id JOIN questions_versions qv ON qv.id = aqa.question_version_id WHERE a.user_id = $1 AND qv.question_id = $2", student_id, id_question).Scan(&id_atempt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с параметрами видимо"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT aqa.question_version_id FROM atempts a JOIN tests t ON t.id = a.test_id JOIN atempts_questions_answers aqa ON a.id = aqa.atempt_id JOIN questions_versions qv ON aqa.question_version_id = qv.id  WHERE t.deleted = false AND a.id = $1 AND qv.question_id = $2", id_atempt, id_question).Scan(&question_version_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с параметрами видимо"})
	}
	_, err = db.Pool.Exec(context.Background(), "UPDATE  atempts_questions_answers SET answer_id = -1 WHERE atempt_id = $1 AND question_version_id = $2", id_atempt, question_version_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с чем то..."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
