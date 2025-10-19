package models

type User struct {
	ID       int64  `json:"id" example:"1"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"mypassword123"`
}

type NewUser struct {
	Email    string
	Password string
}
