package department

import "github.com/jmoiron/sqlx"

type Repository interface {
	FindAll() ([]Department, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Department, error) {
	var departments []Department

	query := `
		SELECT id, name
		FROM departments
		ORDER BY name ASC
	`

	err := r.db.Select(&departments, query)
	return departments, err
}
