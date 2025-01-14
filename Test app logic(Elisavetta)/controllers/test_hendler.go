package controllers

import (
	"context"
	"log"
	"strconv"

	"core_mod/db"
	"core_mod/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func DeletedQuestionFromTest(c *fiber.Ctx) error {
	var num int
	id_t, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id теста не распарсилось "})
	}
	q_id, err := strconv.Atoi(c.Params("q_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id вопроса не распарсилось "})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM atempts WHERE test_id = $1", id_t).Scan(&num)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}
	if num > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Немогу удалить вопрос, тест уже проходили"})
	}
	commandTag, err := db.Pool.Exec(c.Context(), "DELETE FROM public.tests_questions WHERE test_id = $1 AND question_id = $2", id_t, q_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func AddQuestionToTest(c *fiber.Ctx) error {
	var num int
	id_t, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id теста не распарсилось "})
	}
	q_id, err := strconv.Atoi(c.Params("q_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id вопроса не распарсилось "})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM atempts WHERE test_id = $1", id_t).Scan(&num)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}
	if num > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Немогу вставить вопрос тест уже проходили"})
	}
	commandTag, err := db.Pool.Exec(c.Context(), "INSERT INTO tests_questions( test_id, question_id) VALUES ($1, $2)", id_t, q_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func ChangeQuestionOrderInTest(c *fiber.Ctx) error {
	var num int
	var q_ids []int
	id_t, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id теста не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM atempts WHERE test_id = $1", id_t).Scan(&num)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет вопроса с таким id"})
	}
	if num > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Немогу помеянять порядок вопросов,тест уже проходили"})
	}
	if err := c.BodyParser(&q_ids); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустой булевый статус "})
	}
	tx, err := db.Pool.BeginTx(c.Context(), pgx.TxOptions{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось начать транзакцию ": err})
	}
	for i, id := range q_ids {
		commandTag, err := tx.Exec(context.Background(), "UPDATE tests_questions SET q_order = $1 WHERE test_id =$2 AND question_id = $3", i, id_t, id)
		if err != nil {
			tx.Rollback(c.Context())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет теста с таким id"})
		}
		if commandTag.RowsAffected() != 1 {
			tx.Rollback(c.Context())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
		}
	}
	err = tx.Commit(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Не удалось завершить транзакцию ": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func GetUsersPassedTest(c *fiber.Ctx) error {
	var response []int
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id теста не распарсилось "})
	}
	rows, err := db.Pool.Query(context.Background(), "SELECT  user_id from atempts where test_id = $1 AND active = false ", id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u_id int
		err := rows.Scan(&u_id)
		if err != nil {
			log.Println("Нераспарсился id в методе GetUsersPassedTest: ", err)
		}
		response = append(response, u_id)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func GetUserMarksTest(c *fiber.Ctx) error {
	var response []models.UserMark
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id теста не распарсилось "})
	}
	rows, err := db.Pool.Query(context.Background(), "SELECT a.user_id,  SUM(CASE WHEN aqa.answer_id = qv.corect_answer_id THEN 1 ELSE 0 END)*100/COUNT(answer_id) AS mark FROM atempts_questions_answers aqa JOIN questions_versions qv ON aqa.question_version_id = qv.id JOIN atempts a ON a.id = aqa.atempt_id  WHERE a.test_id = $1 GROUP BY a.user_id, a.id", id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserMark
		err := rows.Scan(&user.Id, &user.Mark)
		if err != nil {
			log.Println("Нераспарсился user из базы в методе GetUserMarksTest: ", err)
		}
		response = append(response, user)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func GetUserAnswersTest(c *fiber.Ctx) error {
	// TODO нужно переделать
	var response []models.UserAnswers
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id теста не распарсилось "})
	}
	rows, err := db.Pool.Query(context.Background(), "SELECT at.user_id, qv.text_q, an.title FROM atempts at JOIN atempts_questions_answers aqa ON at.id = aqa.atempt_id JOIN questions_versions qv ON qv.id = aqa.question_version_id JOIN answers an ON an.id = aqa.answer_id WHERE at.test_id = $1 ORDER BY at.user_id", id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()
	var userAnswers models.UserAnswers

	for rows.Next() {
		var answer models.UserAnswer
		var curUserId int

		err := rows.Scan(&curUserId, &answer.QuestionText, &answer.AnsverText)
		if err != nil {
			log.Println("Нераспарсился user из базы в методе GetUserMarksTest: ", err)
		}
		if userAnswers.UserID == 0 {
			userAnswers.UserID = curUserId
			userAnswers.Answers = append(userAnswers.Answers, answer)
		} else if userAnswers.UserID == curUserId {
			userAnswers.Answers = append(userAnswers.Answers, answer)
		} else {
			response = append(response, userAnswers)
			var tmp models.UserAnswers
			tmp.UserID = curUserId
			tmp.Answers = append(tmp.Answers, answer)
			userAnswers = tmp
		}
	}
	response = append(response, userAnswers)

	return c.Status(fiber.StatusOK).JSON(response)
}
