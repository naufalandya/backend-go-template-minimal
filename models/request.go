package model

// Generated structs from C:\Andya\Go\standard_template\mock\request.txt

type Userrequest struct {
  UserId string `json:"user_id" validate:"required,min=1,max=255"`
}