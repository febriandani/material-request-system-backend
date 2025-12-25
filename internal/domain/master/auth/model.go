package auth

type Authentication struct {
	UserID   int64  `db:"user_id"`
	Password string `db:"password"`
}
