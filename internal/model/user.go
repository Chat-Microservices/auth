package model

import "time"

type User struct {
	ID        int64
	Detail    Detail
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	ID    int64
	Name  string
	Email string
}

type Detail struct {
	Name  string
	Email string
	Role  int
}
