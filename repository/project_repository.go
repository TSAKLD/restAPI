package repository

import (
	"context"
	"database/sql"
	"errors"
	"restAPI/entity"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(database *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		db: database,
	}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	q := "INSERT INTO projects(name, user_id, created_at) VALUES($1, $2, $3) RETURNING id"

	tx, err := r.db.Begin()
	if err != nil {
		return entity.Project{}, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, q, project.Name, project.UserID, project.CreatedAt).Scan(&project.ID)
	if err != nil {
		return entity.Project{}, err
	}

	err = r.addProjectMember(ctx, tx, project.ID, project.UserID)
	if err != nil {
		return entity.Project{}, err
	}

	return project, tx.Commit()
}

func (r *ProjectRepository) UserProjects(ctx context.Context, userID int64) (projects []entity.Project, err error) {
	q := "SELECT p.id, p.name, p.user_id, p.created_at FROM projects p JOIN projects_users pu ON pu.project_id = p.id WHERE pu.user_id = $1"

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

func (r *ProjectRepository) ProjectByID(ctx context.Context, id int64) (p entity.Project, err error) {
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

func (r *ProjectRepository) DeleteProject(ctx context.Context, projectID int64) error {
	q := "DELETE FROM projects WHERE id = $1"

	_, err := r.db.ExecContext(ctx, q, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) AddProjectMember(ctx context.Context, projectID int64, userID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.addProjectMember(ctx, tx, projectID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ProjectRepository) addProjectMember(ctx context.Context, tx *sql.Tx, projectID int64, userID int64) error {
	q := "INSERT INTO projects_users(project_id, user_id) VALUES ($1, $2)"

	_, err := tx.ExecContext(ctx, q, projectID, userID)
	if err != nil {
		return err
	}

	return nil
}
