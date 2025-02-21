package models

type User struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Created   int    `json:"created"`
	Updated   int    `json:"updated"`
	Active    int    `json:"active"`
}

type LoginParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
