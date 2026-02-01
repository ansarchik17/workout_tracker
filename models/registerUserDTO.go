package models

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
