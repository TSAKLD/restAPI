package repository

import (
	"context"
	"database/sql"
	"errors"
	"restAPI/entity"
)

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
