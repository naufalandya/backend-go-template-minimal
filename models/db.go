package model

// Generated structs from C:\Andya\Go\modular monolith\mock\db.txt

type User struct {
  Id int32 `json:"id" validate:"required,min=1"`
  Name string `json:"name" validate:"required,min=1,max=255"`
  Email string `json:"email" validate:"required,min=1,max=255"`
}