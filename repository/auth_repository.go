package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"restAPI/entity"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) UserByEmailAndPassword(ctx context.Context, email string, password string) (u entity.User, err error) {
	q := "SELECT id, name, email, created_at, is_verified FROM users WHERE email = $1 AND password = $2"

	err = r.db.QueryRowContext(ctx, q, email, password).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.IsVerified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *AuthRepository) ProjectUsers(ctx context.Context, projectID int64) (users []entity.User, err error) {
	q := `SELECT u.id, u.name, u.email, u.created_at, u.is_verified
	FROM users u
	    JOIN projects_users pu ON pu.user_id = u.id
	WHERE pu.project_id = $1`

	rows, err := r.db.QueryContext(ctx, q, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.IsVerified)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, sessionID uuid.UUID, userID int64, createdAt time.Time) error {
	q := "INSERT INTO sessions(id, user_id, created_at) VALUES($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, q, sessionID, userID, createdAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) UserBySessionID(ctx context.Context, sessionID string) (u entity.User, err error) {
	q := "SELECT u.id, u.email, u.name, u.created_at, u.is_verified FROM users u JOIN sessions s ON u.id = s.user_id WHERE s.id = $1"

	err = r.db.QueryRowContext(ctx, q, sessionID).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.IsVerified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *AuthRepository) SaveVerificationCode(ctx context.Context, code string, userID int64) error {
	q := "INSERT INTO verification_codes(code, user_id) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, q, code, userID)
	return err
}

func (r *AuthRepository) VerifyUser(ctx context.Context, code string) error {
	q := "SELECT user_id FROM verification_codes WHERE code = $1 "

	var id int64

	err := r.db.QueryRowContext(ctx, q, code).Scan(&id)
	if err != nil {
		return err
	}

	q = "UPDATE users SET is_verified = TRUE WHERE id = $1"

	_, err = r.db.ExecContext(ctx, q, id)
	return err
}
