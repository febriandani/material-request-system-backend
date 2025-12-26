package user

type User struct {
	ID           int64  `db:"id" json:"id"`
	FullName     string `db:"full_name" json:"fullName"`
	Role         string `db:"role" json:"role"`
	DepartmentID int64  `db:"department_id" json:"departmentId"`
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	Phone        string `db:"phone" json:"phone"`
}
