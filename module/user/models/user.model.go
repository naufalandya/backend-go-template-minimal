package models

type User struct {
	Id    int32  `json:"id" validate:"required,min=1"`
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,min=1,max=255"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
	FullName string `json:"full_name" validate:"required,min=1,max=255"`
	Username string `json:"username" validate:"required,min=1,max=255"`
}
