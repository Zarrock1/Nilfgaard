package controllers

import (
	"context"
	"core_mod/db"
	"core_mod/models"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func QuestionsHendler(c *fiber.Ctx) error {
	var response []models.Questions
	rows, err := db.Pool.Query(context.Background(), "WITH RankedQuestions AS (SELECT qv.question_id, qv.title, qv.version, ROW_NUMBER() OVER (PARTITION BY qv.question_id ORDER BY qv.version DESC) AS rn FROM questions_versions qv  ) SELECT  q.avtor_id, rq.title, rq.version FROM RankedQuestions rq JOIN questions q ON rq.question_id = q.id WHERE rq.rn = 1 AND q.deleted = FALSE")
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var questions models.Questions
		err := rows.Scan(&questions.AvtorID, &questions.Name, &questions.Vesion)
		if err != nil {
			log.Println("Нераспарсился questions в методе : QuestionsHendler", err)
		}
		response = append(response, questions)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func QuestionHendler(c *fiber.Ctx) error {
	var response models.Question
	var statusqestion bool
	var version_id int

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Id не распарсилось "})
	}
	version, err := strconv.Atoi(c.Params("v_id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM questions WHERE id = $1", id).Scan(&statusqestion)

	if err != nil || statusqestion {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT qv.title, qv.text_q, a.number, qv.id FROM questions_versions qv  JOIN answers a ON qv.corect_answer_id = a.id WHERE qv.question_id = $1 AND qv.version = $2", id, version).Scan(&response.Title, &response.Text, &response.CorrectAnsver, &version_id)
	response.Id = id
	if err != nil {
		log.Println("Database error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}
	rows, err := db.Pool.Query(context.Background(), "SELECT title FROM answers WHERE question_version_id = $1 ORDER BY number", version_id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var answer string
		err := rows.Scan(&answer)
		if err != nil {
			log.Println("Нераспарсился answer в методе : QuestionHendler", err)
		}
		response.Ansver = append(response.Ansver, answer)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func QuestionUpdate(c *fiber.Ctx) error {
	var commandTag pgconn.CommandTag
	var request models.Question
	var statusqestion bool
	var version int
	var response_id int
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое имя вопроса и не пустое описание"})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM questions WHERE id = $1", id).Scan(&statusqestion)
	if err != nil || statusqestion {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}
	tx, _ := db.Pool.BeginTx(c.Context(), pgx.TxOptions{})
	err = tx.QueryRow(context.Background(), "SELECT  version FROM questions_versions  WHERE question_id = $1 ORDER BY version DESC LIMIT 1", id).Scan(&version)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}
	err = tx.QueryRow(c.Context(), "INSERT INTO questions_versions (question_id, title, text_q, version, corect_answer_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", id, request.Title, request.Text, version, request.CorrectAnsver).Scan(&response_id)
	if err != nil {
		tx.Rollback(c.Context())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}

	strsql := "Insert into answers (title, number, question_version_id) values"
	for i, value := range request.Ansver {
		strsql += fmt.Sprintf("( '%s', %d, %d),", value, i, response_id) //?
	}
	strsql = strsql[:len(strsql)-1]

	commandTag, err = db.Pool.Exec(context.Background(), strsql)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() == 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	err = tx.Commit(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось завершить транзакцию ": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Id": response_id})

}
func QuestionCreate(c *fiber.Ctx) error {
	avtor_id := c.Locals("user_id")
	var commandTag pgconn.CommandTag
	var response_id int
	tx, err := db.Pool.BeginTx(c.Context(), pgx.TxOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось начать транзакцию ": err})
	}
	var request models.Question
	var id int
	if err := c.BodyParser(&request); err != nil || request.Title == "" || request.CorrectAnsver == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое имя вопроса и не пустое описание"})
	}
	err = tx.QueryRow(c.Context(), "INSERT INTO questions (avtor_id, deleted) VALUES ($1, false) RETURNING id", avtor_id).Scan(&id)
	if err != nil {
		tx.Rollback(c.Context())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	err = tx.QueryRow(c.Context(), "INSERT INTO questions_versions (question_id, title, text_q, version, corect_answer_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", id, request.Title, request.Text, 1, request.CorrectAnsver).Scan(&response_id)
	if err != nil {
		tx.Rollback(c.Context())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}

	strsql := "Insert into answers (title, number, question_version_id) values"
	for i, value := range request.Ansver {
		strsql += fmt.Sprintf("( '%s', %d, %d),", value, i, response_id)
	}
	strsql = strsql[:len(strsql)-1]

	commandTag, err = db.Pool.Exec(context.Background(), strsql)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() == 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	err = tx.Commit(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось завершить транзакцию ": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Id": id})
}
func QuestionsDelete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	var num int
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM tests_questions WHERE question_id = $1", id).Scan(&num)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}
	if num > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Немогу удалить вопрос, связан с тестом"})
	}
	commandTag, err := db.Pool.Exec(c.Context(), "UPDATE questions SET deleted = TRUE WHERE id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
