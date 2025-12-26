package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByUsername(username string) (*Authentication, int64, string, error)
	SaveRefreshToken(userID int64, tokenHash string, expiresAt time.Time) error
	FindValidRefreshToken(ctx context.Context, userID int64, tokenHash string) (bool, error)
	GetUserRole(ctx context.Context, userID int64) (string, error)
	RevokeRefreshToken(ctx context.Context, userID int64, tokenHash string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUsername(username string) (*Authentication, int64, string, error) {
	query := `
		SELECT a.user_id, a.password, u.username, u.full_name, u.email, u.phone, u.role
		FROM master.authentications a
		JOIN master.users u ON u.id = a.user_id
		WHERE u.username = $1
	`

	var auth Authentication
	var role string

	err := r.db.QueryRowx(query, username).
		Scan(&auth.UserID, &auth.Password, &auth.Username, &auth.FullName, &auth.Email, &auth.Phone, &role)

	if err != nil {
		return nil, 0, "", errors.New("user not found")
	}

	return &auth, auth.UserID, role, nil
}

func (r *repository) SaveRefreshToken(
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) error {

	query := `
		INSERT INTO master.refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, userID, tokenHash, expiresAt)
	return err
}

func (r *repository) FindValidRefreshToken(
	ctx context.Context,
	userID int64,
	tokenHash string,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM master.refresh_tokens
			WHERE user_id = $1
			  AND token_hash = $2
			  AND revoked = false
			  AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userID, tokenHash)
	return exists, err
}

func (r *repository) GetUserRole(
	ctx context.Context,
	userID int64,
) (string, error) {

	var role string
	err := r.db.GetContext(
		ctx,
		&role,
		`SELECT role FROM master.users WHERE id = $1`,
		userID,
	)
	return role, err
}

func (r *repository) RevokeRefreshToken(ctx context.Context, userID int64, tokenHash string) error {
	query := `
		UPDATE master.refresh_tokens
		SET revoked = true
		WHERE user_id = $1
		  AND token_hash = $2
		  AND revoked = false
	`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash)
	return err
}
