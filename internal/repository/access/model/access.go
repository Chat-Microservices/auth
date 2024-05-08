package model

type Access struct {
	ID     int64  `db:"id"`
	RoleId int64  `db:"role_id"`
	Path   string `db:"path"`
}
