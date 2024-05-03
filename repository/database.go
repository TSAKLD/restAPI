package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"restAPI/entity"
)

type Repository struct {
	db *sql.DB
}

func New(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) CreateUser(u entity.User) (entity.User, error) {
	q := "INSERT INTO users(name, password, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	err := r.db.QueryRow(q, u.Name, u.Password, u.Email, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return entity.User{}, err
	}

	return u, nil
}

func (r *Repository) DeleteUser(id int64) error {
	q := "DELETE FROM users WHERE id = $1"

	_, err := r.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserByID(id int64) (u entity.User, err error) {
	q := "SELECT id, name, email, created_at FROM users WHERE id = $1"

	err = r.db.QueryRow(q, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *Repository) UserByEmail(email string) (u entity.User, err error) {
	q := "SELECT id, name, email, created_at FROM users WHERE email = $1"

	err = r.db.QueryRow(q, email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *Repository) Users() (users []entity.User, err error) {
	q := "SELECT id, name, email, created_at FROM users"

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) UserByEmailAndPassword(email string, password string) (u entity.User, err error) {
	q := "SELECT id, name, email, created_at FROM users WHERE email = $1 AND password = $2"

	err = r.db.QueryRow(q, email, password).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *Repository) CreateSession(sessionID uuid.UUID, userID int64) error {
	q := "INSERT INTO sessions(session_id, user_id) VALUES($1, $2)"

	_, err := r.db.Exec(q, sessionID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserBySessionID(sessionID string) (u entity.User, err error) {
	q := "select u.id, u.email, u.name, u.created_at from users u JOIN sessions s ON u.id = s.user_id WHERE s.session_id = $1"

	err = r.db.QueryRow(q, sessionID).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}
