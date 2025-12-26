package user

import "github.com/jmoiron/sqlx"

type Repository interface {
	FindAll() ([]User, error)
	FindApprovers() ([]Approver, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]User, error) {
	var users []User

	query := `
		SELECT 
		id,
		full_name,
		role,
		department_id,
		username,
		email,
		phone
		FROM master.users ORDER BY full_name ASC
	`

	err := r.db.Select(&users, query)

	return users, err
}

func (r *repository) FindApprovers() ([]Approver, error) {
	var approvers []Approver

	query := `
		SELECT
		u.id,
		CONCAT(u.full_name, ' - Head of ', d.name) AS approver_name,
		u.role,
		u.department_id,
		u.username,
		u.email,
		u.phone
		FROM master.users u
		INNER JOIN master.departments d on d.id = u.department_id
		WHERE u.role = 'approver'
		ORDER BY u.full_name ASC
	`

	err := r.db.Select(&approvers, query)

	return approvers, err
}
