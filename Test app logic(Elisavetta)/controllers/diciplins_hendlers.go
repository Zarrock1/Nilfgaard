package controllers

import (
	"context"

	//"fmt"
	"log"
	"strconv"

	"core_mod/db"
	"core_mod/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

func DisciplinsHendler(c *fiber.Ctx) error {
	var response []models.Disciplin
	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, discription FROM disciplines WHERE deleted= false")
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close() // ?
	for rows.Next() {
		var descriplin models.Disciplin
		err := rows.Scan(&descriplin.Id, &descriplin.DisciplinName, &descriplin.Discription)
		if err != nil {
			log.Println("Нераспарсился disciplines в методе : DisciplinsHendler", err)
		}
		response = append(response, descriplin)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func DisciplinHendler(c *fiber.Ctx) error {
	var response models.DisciplinP
	var statusdisciplin bool
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&statusdisciplin)

	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT title, discription, prepod_id FROM disciplines WHERE id = $1", id).Scan(&response.DisciplinName, &response.Discription, &response.PrepodId)

	if err != nil {
		log.Println("Database error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func DisciplinUpdate(c *fiber.Ctx) error {
	var request models.Disciplin
	var isDisciplineDeleted bool
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Id не распарсилось "})
	}
	if err := c.BodyParser(&request); err != nil || (request.DisciplinName == "" && request.Discription == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое disciplinename или discription"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&isDisciplineDeleted)
	if err != nil || isDisciplineDeleted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}

	request.Id = id
	var commandTag pgconn.CommandTag
	if request.DisciplinName != "" && request.Discription != "" {
		commandTag, err = db.Pool.Exec(context.Background(), "UPDATE disciplines SET title= $1, discription =$2 WHERE id=$3 RETURNING title, id, discription", request.DisciplinName, request.Discription, request.Id)
	} else if request.Discription == "" {
		commandTag, err = db.Pool.Exec(context.Background(), "UPDATE disciplines SET title=$1 WHERE id=$2 RETURNING title", request.DisciplinName, request.Id)
	} else {
		commandTag, err = db.Pool.Exec(context.Background(), "UPDATE disciplines SET discription=$1 WHERE id=$2 RETURNING title", request.Discription, request.Id)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func DisciplinTests(c *fiber.Ctx) error {
	var response []models.DisciplinTest
	var statusdisciplin bool
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /disciplins/id где id это цифорка"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	rows, err := db.Pool.Query(context.Background(), "Select id, title from tests where discipline_id=$1 AND deleted = false", id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close() // ?

	for rows.Next() {
		var test models.DisciplinTest
		err := rows.Scan(&test.TestId, &test.TestName)
		if err != nil {
			log.Println("Нераспарсился test в методе DisciplinTests: ", err)
		}
		response = append(response, test)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func DisciplinTestStatus(c *fiber.Ctx) error {
	var response models.TestStatus
	var statusdisciplin bool
	var teststatus bool
	t_id, err := strconv.Atoi(c.Params("t_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT disciplines.deleted FROM tests JOIN disciplines ON tests.discipline_id = disciplines.id WHERE tests.id =$1 AND disciplines.id =$2", t_id, id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM tests WHERE id=$1", t_id).Scan(&teststatus)
	if err != nil || teststatus {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет теста с таким id"})
	}
	err = db.Pool.QueryRow(context.Background(), "Select active from tests  where id=$1", t_id).Scan(&response.Teststatus)
	if err != nil {
		log.Println("Database error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func DisciplinTestStatusUpdate(c *fiber.Ctx) error {
	var request models.TestStatus
	var statusdisciplin bool
	var teststatusdeleted bool
	t_id, err := strconv.Atoi(c.Params("t_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустой булевый статус "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT disciplines.deleted FROM tests JOIN disciplines ON tests.discipline_id = disciplines.id WHERE tests.id =$1", t_id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM tests WHERE id=$1", t_id).Scan(&teststatusdeleted)
	if err != nil || teststatusdeleted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет теста с таким id"})
	}

	commandTag, err := db.Pool.Exec(context.Background(), "UPDATE tests SET active = $1 WHERE id =$2", request.Teststatus, t_id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет теста с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func DisciplinTestCreaite(c *fiber.Ctx) error {
	var request models.DisciplinTest
	var test_id int
	var statusdisciplin bool
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	if err := c.BodyParser(&request); err != nil || request.TestName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое имя теста "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	err = db.Pool.QueryRow(context.Background(), "INSERT INTO tests (title, active, discipline_id) VALUES ($1, false, $2) RETURNING id", request.TestName, id).Scan(&test_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет дисциплины с таким id"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ID": test_id})
}
func DisciplinTestDelete(c *fiber.Ctx) error {
	var statusdisciplin bool
	t_id, err := strconv.Atoi(c.Params("t_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT disciplines.deleted FROM tests JOIN disciplines ON tests.discipline_id = disciplines.id WHERE tests.id =$1 AND disciplines.id =$2", t_id, id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	commandTag, err := db.Pool.Exec(context.Background(), "Update tests Set deleted = true Where id = $1", t_id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func DisciplinStudents(c *fiber.Ctx) error {
	var response []models.DisciplinStudent
	var statusdisciplin bool
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /disciplins/id/students где id это цифорка"})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}

	rows, err := db.Pool.Query(context.Background(), "SELECT  user_id from users_disciplines where discipline_id =$1", id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var student models.DisciplinStudent
		err := rows.Scan(&student.StudentId)
		if err != nil {
			log.Println("Нераспарсился student в методе DisciplinStudents: ", err)
		}
		response = append(response, student)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func DisciplinStudentAdd(c *fiber.Ctx) error {
	var statusdisciplin bool
	var commandTag pgconn.CommandTag

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /disciplins/id/students/students/:s_id где id это цифорка дисциплины"})
	}
	s_id, err := strconv.Atoi(c.Params("s_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /disciplins/id/students/students/:s_id где s_id это цифорка студента"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	commandTag, err = db.Pool.Exec(context.Background(), "INSERT INTO users_disciplines (user_id, discipline_id) VALUES ($1,$2)", s_id, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет дисциплины с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func DisciplinStudentDelete(c *fiber.Ctx) error {
	var statusdisciplin bool
	var commandTag pgconn.CommandTag

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /disciplins/id/students/students/:s_id где id это цифорка дисциплины"})
	}
	s_id, err := strconv.Atoi(c.Params("s_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /disciplins/id/students/students/:s_id где s_id это цифорка студента"})
	}
	err = db.Pool.QueryRow(context.Background(), "SELECT deleted FROM disciplines WHERE id = $1", id).Scan(&statusdisciplin)
	if err != nil || statusdisciplin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Нет дисциплины с таким id"})
	}
	commandTag, err = db.Pool.Exec(context.Background(), "DELETE FROM users_disciplines WHERE user_id = $1 AND discipline_id = $2", s_id, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет дисциплины с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func DisciplinCreate(c *fiber.Ctx) error {
	var request models.DisciplinP
	var id int
	if err := c.BodyParser(&request); err != nil || request.DisciplinName == "" || request.Discription == "" || request.PrepodId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое имя десциплины и не пустое описание и не пустой id препода "})
	}
	err := db.Pool.QueryRow(context.Background(), "INSERT INTO disciplines (title, discription, prepod_id, deleted) VALUES ($1, $2, $3, false) RETURNING id", request.DisciplinName, request.Discription, request.PrepodId).Scan(&id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Упс ошибочка вышла"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Id": id})
}
func DisciplinDeleted(c *fiber.Ctx) error {
	var commandTag pgconn.CommandTag
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	commandTag, err = db.Pool.Exec(context.Background(), "UPDATE disciplines SET deleted = true WHERE id =$1", id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет дисциплины с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
