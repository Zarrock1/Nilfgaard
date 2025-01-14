package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"core_mod/db"
	"core_mod/models"

	"github.com/gofiber/fiber/v2"
)

func UserHandler(c *fiber.Ctx) error {

	var response models.User

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}

	response.Id = id

	err = db.Pool.QueryRow(context.Background(), "SELECT name FROM users WHERE id=$1", id).Scan(&response.Username)
	if err != nil {
		log.Println("Database error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func UsersHandler(c *fiber.Ctx) error {
	var response []models.User

	rows, err := db.Pool.Query(context.Background(), "SELECT id, name FROM users")
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close() // ?
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			log.Println("Нераспарсился user в методе UsersHendler: ", err)
		}
		response = append(response, user)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func UserUpdate(c *fiber.Ctx) error {
	// Заменяет имя по id

	var request models.User
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /users/id где id это цифорка"})
	}

	if err := c.BodyParser(&request); err != nil || request.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое username "})
	}

	request.Id = id
	commandTag, err := db.Pool.Exec(context.Background(), "UPDATE users SET name=$1 WHERE id=$2 RETURNING name", request.Username, request.Id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(request)
}
func UserTests(c *fiber.Ctx) error {
	var response []models.UserDisciplinTests
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	rows, err := db.Pool.Query(context.Background(), "SELECT d.title, t.title, SUM(CASE WHEN aqa.answer = qv.corect_answer_id THEN 1 ELSE 0 END)*100/COUNT(answer) AS mark FROM atemps_questions_answers aqa JOIN questions_versions qv ON aqa.question_version_id = qv.id JOIN atemps a ON a.id = aqa.atempt_id JOIN tests t ON t.id = a.test_id JOIN disciplines d ON d.id = t.discipline_id WHERE a.user_id = $1 GROUP BY a.test_id, t.title, d.title, a.id ", id)
	if err != nil {
		log.Println("Database error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}
	defer rows.Close()
	var curUDT models.UserDisciplinTests
	var temp models.UserTest
	var temp_dname string
	var temp_name string
	var temp_mark int
	// Здесь сканируем первый результат перед циклом
	if rows.Next() {

		err = rows.Scan(&temp_dname, &temp_name, &temp_mark)
		if err != nil {
			log.Println("Нераспарсился первый результат в методе UsersTests: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Ошибка при первом сканировании"})
		}

		curUDT.DisciplinName = temp_dname
		temp.Name = temp_name
		temp.Mark = temp_mark
		curUDT.Tests = append(curUDT.Tests, temp)
	}

	for rows.Next() {

		// var temp models.Test
		//var temp_dname string

		err = rows.Scan(&temp_dname, &temp_name, &temp_mark)
		if err != nil {
			log.Println("Нераспарсился test в методе UsersTests: ", err)
		}
		temp.Name = temp_name
		temp.Mark = temp_mark

		if curUDT.DisciplinName != temp_dname {
			response = append(response, curUDT)

			curUDT.DisciplinName = temp_dname
			curUDT.Tests = make([]models.UserTest, 0)
		}
		curUDT.Tests = append(curUDT.Tests, temp)
	}
	response = append(response, curUDT)

	return c.Status(fiber.StatusOK).JSON(response)
}
func UserRoles(c *fiber.Ctx) error {
	var response []models.Role

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}

	rows, err := db.Pool.Query(context.Background(), "SELECT roles.id, roles.name FROM users_roles JOIN roles ON users_roles.role_id = roles.id WHERE users_roles.user_id=$1", id)
	if err != nil {
		log.Println("Database error: ", err)
	}
	defer rows.Close() // ?

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.Id, &role.RoleName)
		if err != nil {
			log.Println("Нераспарсился role в методе UsersRoles: ", err)
		}
		response = append(response, role)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func UserUpdateRoles(c *fiber.Ctx) error {
	var request []models.Role
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /users/id/role где id это цифорка"})

	}
	if err := c.BodyParser(&request); err != nil || (request == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать что-то"})
	}
	commandTag, err := db.Pool.Exec(context.Background(), "Delete from users_roles where user_id = $1", id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() == 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	strsql := "Insert into users_roles (user_id, role_id) values"
	for _, value := range request {
		strsql += fmt.Sprintf("( %d, %d),", id, value.Id)
	}
	strsql = strsql[:len(strsql)-1]

	commandTag, err = db.Pool.Exec(context.Background(), strsql)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() == 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Всё хорошо, всё под контролем"})
}
func UserStatus(c *fiber.Ctx) error {
	var response models.UserStatus

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}

	err = db.Pool.QueryRow(context.Background(), "SELECT blocked FROM users WHERE id=$1", id).Scan(&response.Status)
	if err != nil {
		log.Println("Database error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func UserUpdateStatus(c *fiber.Ctx) error {
	var request models.UserStatus
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /users/id/status где id это цифорка"})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустой булевый статус "})
	}

	commandTag, err := db.Pool.Exec(context.Background(), "UPDATE users SET blocked=$1 WHERE id=$2 RETURNING blocked", request.Status, id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "Наверное нет пользователя с таким id"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Упс... Ошибочка вышла (лажа в базе наверное)"})
	}

	return c.Status(fiber.StatusOK).JSON(request)
}
