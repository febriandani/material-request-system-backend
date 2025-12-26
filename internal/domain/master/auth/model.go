package auth

type Authentication struct {
	UserID   int64  `db:"user_id"`
	FullName string `db:"full_name"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}
