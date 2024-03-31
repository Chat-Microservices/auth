package modelRepo

import "time"

type User struct {
	ID        int64     `db:"id"`
	Detail    Detail    `db:""`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdateUserRequest struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

type Detail struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  int    `db:"role"`
}
