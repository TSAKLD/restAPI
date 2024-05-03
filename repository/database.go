package repository

import (
	"database/sql"
	"errors"
	"restAPI/entity"
	"restAPI/service"
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
			return entity.User{}, service.ErrNotFound
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
			return entity.User{}, service.ErrNotFound
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
