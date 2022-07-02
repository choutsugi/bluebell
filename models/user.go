package models

type User struct {
	Uid      int64  `db:"uid"`
	Username string `db:"username"`
	Password string `db:"password"`
}
