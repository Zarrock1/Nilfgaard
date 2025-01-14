package models

type User struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

type Usreresponse struct {
	Disciplins string `json:"disciplins"`
	Group      string `json:"group"`
	Course     int    `json:"course"`
	Tests_user UserTest
}

type UserTest struct {
	Name string `json:"name"`
	Mark int    `json:"mark"`
}

type Role struct {
	Id       int    `json:"id"`
	RoleName string `json:"rolename"`
}
type UserStatus struct {
	Status bool `json:"blocked"`
}
type UserDisciplinTests struct {
	DisciplinName string     `json:"disciplinname"`
	Tests         []UserTest `json:"tests"`
}

type Disciplin struct {
	Id            int    `json:"id"`
	DisciplinName string `json:"name"`
	Discription   string `json:"discription"`
}

type DisciplinP struct {
	PrepodId      int    `json:"prepod_id"`
	DisciplinName string `json:"name"`
	Discription   string `json:"discription"`
}

type DisciplinTest struct {
	TestName string `json:"name"`
	TestId   int    `json:"id"`
}

type TestStatus struct {
	Teststatus bool `json:"active"`
}
type DisciplinStudent struct {
	StudentId int `json:"student"`
}
type Questions struct {
	Name    string `json:"name"`
	Vesion  int    `json:"version"`
	AvtorID int    `json:"avtor_id"`
}
type Question struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	Text          string   `json:"text"`
	Ansver        []string `json:"ansvers"`
	CorrectAnsver int      `json:"coretansver"`
}
type UserMark struct {
	Id   int `json:"id"`
	Mark int `json:"mark"`
}
type Atempt struct {
	Status_active bool         `json:"active"`
	Answers       []UserAnswer `jsons:"answers"`
}
type UserAnswers struct {
	UserID  int          `json:"id"`
	Answers []UserAnswer `jsons:"answers"`
}

type UserAnswer struct {
	QuestionText string `json:"text"`
	AnsverText   string `json:"ansvers"`
}
