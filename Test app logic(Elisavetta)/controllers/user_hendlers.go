package controllers

import (
	"log"
	"strconv"

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
	//TODO: нужно добавить логику работы с базой
	response.Username = "FakeUser"

	return c.Status(fiber.StatusOK).JSON(response)
}

func UsersHandler(c *fiber.Ctx) error {
	var response []models.User

	//TODO: нужно добавить логику работы с базой

	return c.Status(fiber.StatusOK).JSON(response)
}

func UserUpdate(c *fiber.Ctx) error {

	var request models.User
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось правильный запрос долежн быть такой /users/id где id это цифорка"})
	}

	if err := c.BodyParser(&request); err != nil || request.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Message": "invalide body должно содержать не пустое username "})
	}

	request.Id = id

	//TODO: нужно добавить логику работы с базой

	return c.Status(fiber.StatusOK).JSON(request)
}

func UserRoles(c *fiber.Ctx) error {
	var response []models.Role

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Id не распарсилось "})
	}
	log.Println(id)

	//TODO: нужно добавить логику работы с базой

	return c.Status(fiber.StatusOK).JSON(response)
}
