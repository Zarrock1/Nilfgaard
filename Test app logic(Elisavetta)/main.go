package main

import (
	"core_mod/controllers"
	"core_mod/db"
	gwt "core_mod/jwt"
	"core_mod/privelegies"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Подключаемся к базе данных
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		db.CloseDB()
		log.Print("db closed")
	}()

	app := fiber.New()
	app.Use(logger.New())
	app.Static("/", "./public")

	app.Post("/gettoken", gwt.TokenHandler)

	const base_api_url = "/api"
	const base_test_api_url = "/testapi"
	// Тестовое API доступное без токена и проверок  прав доступа
	app.Get(base_test_api_url+"/users", controllers.UsersHandler)
	app.Get(base_test_api_url+"/users/:id", controllers.UserHandler)
	app.Patch(base_test_api_url+"/users/:id/name", controllers.UserUpdate)
	app.Get(base_test_api_url+"/users/:id/tests", controllers.UserTests)
	app.Get(base_test_api_url+"/users/:id/role", controllers.UserRoles)
	app.Patch(base_test_api_url+"/users/:id/role", controllers.UserUpdateRoles)
	app.Get(base_test_api_url+"/users/:id/status", controllers.UserStatus)
	app.Patch(base_test_api_url+"/users/:id/status", controllers.UserUpdateStatus)

	app.Get(base_test_api_url+"/disciplins", controllers.DisciplinsHendler)
	app.Get(base_test_api_url+"/disciplins/:id", controllers.DisciplinHendler)
	app.Patch(base_test_api_url+"/disciplins/:id", controllers.DisciplinUpdate)
	app.Get(base_test_api_url+"/disciplins/:id/tests", controllers.DisciplinTests)
	app.Get(base_test_api_url+"/disciplins/:d_id/tests/:t_id", controllers.DisciplinTestStatus)
	app.Patch(base_test_api_url+"/disciplins/:d_id/tests/:t_id", controllers.DisciplinTestStatusUpdate)
	app.Post(base_test_api_url+"/disciplins/:id/tests", controllers.DisciplinTestCreaite)
	app.Delete(base_test_api_url+"/disciplins/:id/tests/:t_id", controllers.DisciplinTestDelete)
	app.Get(base_test_api_url+"/disciplins/:id/students", controllers.DisciplinStudents)
	app.Post(base_test_api_url+"/disciplins/:id/students/:s_id", controllers.DisciplinStudentAdd)
	app.Delete(base_test_api_url+"/disciplins/:id/students/:s_id", controllers.DisciplinStudentDelete)
	app.Post(base_test_api_url+"/disciplins", controllers.DisciplinCreate)
	app.Delete(base_test_api_url+"/disciplin/:id", controllers.DisciplinDeleted)

	// API с токеном
	app.Get(base_api_url+"/users", gwt.Protected, privelegies.UsersHandler, controllers.UsersHandler)
	app.Get(base_api_url+"/users/:id", gwt.Protected, privelegies.UserHandler, controllers.UserHandler)
	app.Patch(base_api_url+"/users/:id/name", gwt.Protected, privelegies.UserUpdate, controllers.UserUpdate)
	app.Get(base_api_url+"/users/:id/tests", gwt.Protected, privelegies.UserTests, controllers.UserTests)
	app.Get(base_api_url+"/users/:id/role", gwt.Protected, privelegies.UserRoles, controllers.UserRoles)
	app.Patch(base_api_url+"/users/:id/role", gwt.Protected, privelegies.UserUpdateRoles, controllers.UserUpdateRoles)
	app.Get(base_api_url+"/users/:id/status", gwt.Protected, privelegies.UserStatus, controllers.UserStatus)
	app.Patch(base_api_url+"/users/:id/status", gwt.Protected, privelegies.UserUpdateStatus, controllers.UserUpdateStatus)

	app.Get(base_api_url+"/disciplins", gwt.Protected, privelegies.DisciplinsHandler, controllers.DisciplinsHendler)
	app.Get(base_api_url+"/disciplins/:id", gwt.Protected, privelegies.DisciplinHandler, controllers.DisciplinHendler)
	app.Patch(base_api_url+"/disciplins/:id", gwt.Protected, privelegies.DisciplinUpdate, controllers.DisciplinUpdate)
	app.Get(base_api_url+"/disciplins/:id/tests", gwt.Protected, privelegies.DisciplinTests, controllers.DisciplinTests)
	app.Get(base_api_url+"/disciplins/:id/tests/:t_id", gwt.Protected, privelegies.DisciplinTestStatus, controllers.DisciplinTestStatus)
	app.Patch(base_api_url+"/disciplins/:id/tests/:t_id", gwt.Protected, privelegies.DisciplinTestStatusUpdate, controllers.DisciplinTestStatusUpdate)
	app.Post(base_api_url+"/disciplins/:id/tests", gwt.Protected, privelegies.DisciplinTestCreaite, controllers.DisciplinTestCreaite)
	app.Delete(base_api_url+"/disciplins/:id/tests/:t_id", gwt.Protected, privelegies.DisciplinTestDelete, controllers.DisciplinTestDelete)
	app.Get(base_api_url+"/disciplins/:id/students", gwt.Protected, privelegies.DisciplinStudents, controllers.DisciplinStudents)
	app.Post(base_api_url+"/disciplins/:id/students/:s_id", gwt.Protected, privelegies.DisciplinStudentAdd, controllers.DisciplinStudentAdd)
	app.Delete(base_api_url+"/disciplins/:id/students/:s_id", gwt.Protected, privelegies.DisciplinStudentDelete, controllers.DisciplinStudentDelete)
	app.Post(base_api_url+"/disciplins", gwt.Protected, privelegies.DisciplinCreate, controllers.DisciplinCreate)
	app.Delete(base_api_url+"/disciplins/:id", gwt.Protected, privelegies.DisciplinDeleted, controllers.DisciplinDeleted)

	app.Get(base_api_url+"/questions", gwt.Protected, privelegies.QuestionsHendler, controllers.QuestionsHendler)
	app.Get(base_api_url+"/questions/:id/versions/:v_id", gwt.Protected, privelegies.QuestionHendler, controllers.QuestionHendler)
	app.Patch(base_api_url+"/questions/:id", gwt.Protected, privelegies.QuestionUpdate, controllers.QuestionUpdate)
	app.Post(base_api_url+"/questions", gwt.Protected, privelegies.QuestionCreate, controllers.QuestionCreate)
	app.Delete(base_api_url+"/questions/:id", gwt.Protected, privelegies.QuestionsDelete, controllers.QuestionsDelete)

	app.Delete(base_api_url+"/tests/:id/questions/:q_id", gwt.Protected, privelegies.DeletedQuestionFromTest, controllers.DeletedQuestionFromTest)
	app.Post(base_api_url+"/tests/:id/questions/:q_id", gwt.Protected, privelegies.AddQuestionToTest, controllers.AddQuestionToTest)
	app.Post(base_api_url+"/tests/:id/questions", gwt.Protected, privelegies.ChangeQuestionOrderInTest, controllers.ChangeQuestionOrderInTest)
	app.Get(base_api_url+"/tests/:id/users", gwt.Protected, privelegies.GetUsersPassedTest, controllers.GetUsersPassedTest)
	app.Get(base_api_url+"/tests/:id/users/marks", gwt.Protected, privelegies.GetUserMarksTest, controllers.GetUserMarksTest)
	app.Get(base_api_url+"/tests/:id/users/answers", gwt.Protected, privelegies.GetUserAnswersTest, controllers.GetUserAnswersTest)

	app.Post(base_api_url+"/atempts/tests/:t_id", gwt.Protected, privelegies.CreateAtempt, controllers.CreateAtempt)
	app.Patch(base_api_url+"/atempts/:id/questions/:q_id", gwt.Protected, privelegies.UpdateAtempt, controllers.UpdateAtempt)
	app.Patch(base_api_url+"/atempts/:id", gwt.Protected, privelegies.CompleteAtempt, controllers.CompleteAtempt)
	app.Get(base_api_url+"/atempts/tests/:t_id/users/:u_id", gwt.Protected, privelegies.GetAtempts, controllers.GetAtempts)

	app.Get(base_api_url+"/answers/tests/:id", gwt.Protected, privelegies.GetUserAnswersTest, controllers.GetUserAnswersTest)
	app.Patch(base_api_url+"/answers/questions/:q_id", gwt.Protected, privelegies.UapdateAnswer, controllers.UpdateAnswer)
	app.Delete(base_api_url+"/answers/questions/:q_id", gwt.Protected, privelegies.DeleteAnswer, controllers.DeleteAnswer)

	log.Println(app.Listen(":3000"))

}
