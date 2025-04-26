package models

type User struct {
	Id    int32  `json:"id" validate:"required,min=1"`
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,min=1,max=255"`
}

type UserInput struct {
	Name    string   `json:"name" validate:"required,min=1,max=255"`
	Email   string   `json:"email" validate:"required,email"`
	Age     int      `json:"age" validate:"required,gt=0"`
	Hobbies []string `json:"hobbies" validate:"required,dive,required"`
	Profile Profile  `json:"profile" validate:"required"`
}

type Profile struct {
	Bio   string  `json:"bio" validate:"required"`
	Score float64 `json:"score" validate:"required,gt=0"`
}
