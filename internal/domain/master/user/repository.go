package user

import "github.com/jmoiron/sqlx"

type Repository interface {
	FindAll() ([]User, error)
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
