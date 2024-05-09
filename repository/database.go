package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"restAPI/entity"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) CreateUser(ctx context.Context, u entity.User) (entity.User, error) {
	q := "INSERT INTO users(name, password, email, created_at, is_verified) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := r.db.QueryRowContext(ctx, q, u.Name, u.Password, u.Email, u.CreatedAt, u.IsVerified).Scan(&u.ID)
	if err != nil {
		return entity.User{}, err
	}

	return u, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	q := "DELETE FROM users WHERE id = $1"

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserByID(ctx context.Context, id int64) (u entity.User, err error) {
	q := "SELECT id, name, email, created_at, is_verified FROM users WHERE id = $1"

	err = r.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.IsVerified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (u entity.User, err error) {
	q := "SELECT id, name, email, created_at, is_verified FROM users WHERE email = $1"

	err = r.db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.IsVerified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return u, err
	}

	return u, nil
}

func (r *Repository) Users(ctx context.Context) (users []entity.User, err error) {
	q := "SELECT id, name, email, created_at, is_verified FROM users"

	rows, err := r.db.QueryContext(ctx, q)
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

func (r *Repository) UserByEmailAndPassword(ctx context.Context, email string, password string) (u entity.User, err error) {
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

func (r *Repository) CreateSession(ctx context.Context, sessionID uuid.UUID, userID int64, createdAt time.Time) error {
	q := "INSERT INTO sessions(id, user_id, created_at) VALUES($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, q, sessionID, userID, createdAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserBySessionID(ctx context.Context, sessionID string) (u entity.User, err error) {
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

func (r *Repository) CreateProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	q := "INSERT INTO projects(name, user_id, created_at) VALUES($1, $2, $3) RETURNING id"

	err := r.db.QueryRowContext(ctx, q, project.Name, project.UserID, project.CreatedAt).Scan(&project.ID)
	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (r *Repository) UserProjects(ctx context.Context, userID int64) (projects []entity.Project, err error) {
	q := "SELECT id, name, user_id, created_at FROM projects WHERE user_id = $1"

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p entity.Project

		err = rows.Scan(&p.ID, &p.Name, &p.UserID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *Repository) ProjectByID(ctx context.Context, id int64) (p entity.Project, err error) {
	q := "SELECT id, name, user_id, created_at FROM projects WHERE id = $1"

	err = r.db.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.Name, &p.UserID, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Project{}, entity.ErrNotFound
		}

		return p, err
	}

	return p, nil
}

func (r *Repository) DeleteProject(ctx context.Context, projectID int64) error {
	q := "DELETE FROM projects WHERE id = $1"

	_, err := r.db.ExecContext(ctx, q, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	q := "INSERT INTO tasks (name, project_id, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	err := r.db.QueryRowContext(ctx, q, t.Name, t.ProjectID, t.UserID, t.CreatedAt).Scan(&t.ID)
	if err != nil {
		return entity.Task{}, err
	}

	return t, nil
}

func (r *Repository) TaskByID(ctx context.Context, id int64) (t entity.Task, err error) {
	q := "SELECT id, name, project_id, user_id, created_at FROM tasks WHERE id = $1"

	err = r.db.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Name, &t.ProjectID, &t.UserID, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, entity.ErrNotFound
		}

		return t, err
	}

	return t, nil
}

func (r *Repository) SaveVerificationCode(ctx context.Context, code string, userID int64) error {
	q := "INSERT INTO verification_codes(code, user_id) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, q, code, userID)
	return err
}

func (r *Repository) Verify(ctx context.Context, code string) error {
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

func (r *Repository) ProjectTasks(ctx context.Context, projectID int64) (tasks []entity.Task, err error) {
	q := "SELECT id, name, project_id, user_id, created_at FROM tasks WHERE project_id = $1"

	rows, err := r.db.QueryContext(ctx, q, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task entity.Task

		err = rows.Scan(&task.ID, &task.Name, &task.ProjectID, &task.UserID, &task.CreatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) UserTasks(ctx context.Context, userID int64) (tasks []entity.Task, err error) {
	q := "	SELECT id, name, project_id, user_id, created_at FROM tasks WHERE user_id = $1"

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task entity.Task

		err = rows.Scan(&task.ID, &task.Name, &task.ProjectID, &task.UserID, &task.CreatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
