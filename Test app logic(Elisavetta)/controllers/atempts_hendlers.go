package controllers

import (
	"context"
	"core_mod/db"
	"core_mod/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func CreateAtempt(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("t_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	student_id := c.Locals("user_id")
	var commandTag pgconn.CommandTag
	var num int
	var atempt_id int
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM atempts a JOIN tests t ON t.id = a.test_id  where t.active = true AND a.user_id = $1 AND t.id = $2", student_id, id).Scan(&num)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс ошибка рабты с базой"})
	}
	if num > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нельзя либо тест закрыт либо попытки кончились или и "})
	}
	tx, err := db.Pool.BeginTx(c.Context(), pgx.TxOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось начать транзакцию ": err})
	}

	err = tx.QueryRow(c.Context(), "INSERT INTO atempts (user_id, test_id, active) VALUES ($1, $2, TRUE) RETURNING id", student_id, id).Scan(&atempt_id)
	if err != nil {
		tx.Rollback(c.Context())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	commandTag, err = tx.Exec(c.Context(), "INSERT INTO atempts_questions_answers(atempt_id, question_version_id, answer_id)  WITH RankedQuestions AS (SELECT qv.question_id, qv.title, qv.version, qv.id, ROW_NUMBER() OVER (PARTITION BY qv.question_id ORDER BY qv.version DESC) AS rn FROM questions_versions qv  ) SELECT $1, rq.id, -1 FROM RankedQuestions rq JOIN tests_questions tq ON rq.question_id = tq.question_id WHERE rq.rn = 1 AND tq.test_id = $2  ", atempt_id, id)
	if err != nil {
		tx.Rollback(c.Context())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	if commandTag.RowsAffected() < 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	err = tx.Commit(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось завершить транзакцию ": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Id": atempt_id})
}
func UpdateAtempt(c *fiber.Ctx) error {
	var question_version_id int
	var answer_number int
	var answer_id int
	answer_number = -1
	id_atempt, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	id_question, err := strconv.Atoi(c.Params("q_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id question не распарсилось "})
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
func CompleteAtempt(c *fiber.Ctx) error {
	id_atempt, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	commandTag, err := db.Pool.Exec(c.Context(), "UPDATE  atempts SET  active = false WHERE id = $1 ", id_atempt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})

}
func GetAtempts(c *fiber.Ctx) error {
	var response models.Atempt
	var id_atempt int
	id_test, err := strconv.Atoi(c.Params("t_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id попытки не распарсилось "})
	}
	id_user, err := strconv.Atoi(c.Params("u_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id_user не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT active, id FROM atempts Where id = $1 AND user_id = $2", id_test, id_user).Scan(&response.Status_active, &id_atempt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Лажа с параметрами видимо"})
	}
	rows, err := db.Pool.Query(context.Background(), "SELECT qv.text_q, a.title FROM atempts_questions_answers aqa JOIN questions_versions qv ON aqa.question_version_id = qv.id JOIN answers a ON aqa.answer_id = a.id WHERE atempt_id = $1", id_atempt)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var answer models.UserAnswer
		err := rows.Scan(&answer.QuestionText, &answer.AnsverText)
		if err != nil {
			log.Println("Нераспарсился answer в методе : GetAtempts", err)
		}
		response.Answers = append(response.Answers, answer)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
