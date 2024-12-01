package models

type User struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}
type Users struct {
	Usernames []Users
}

type Usreresponse struct {
	Disciplins string `json:"disciplins"`
	Group      string `json:"group"`
	Course     int    `json:"course"`
	Tests_user Tests
}

type Tests struct {
	Name string `json:"tests"`
	Ball int    `json:"ball"`
}

type Role struct {
	RoleName string `json:"rolename"`
}
type UserStatus struct {
	Status bool `json:"rasblokirovan"`
}

type Disciplin struct {
}
