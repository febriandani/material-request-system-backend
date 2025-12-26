package user

type Approver struct {
	ID           int64  `db:"id" json:"userId"`
	ApproverName string `db:"approver_name" json:"approverName"`
	Role         string `db:"role" json:"role"`
	DepartmentID int64  `db:"department_id" json:"departmentId"`
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	Phone        string `db:"phone" json:"phone"`
}
