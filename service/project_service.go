package service

import (
	"context"
	"fmt"
	"restAPI/entity"
	"time"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, t entity.Task) (entity.Task, error)
	TaskByID(ctx context.Context, id int64) (t entity.Task, err error)
	ProjectTasks(ctx context.Context, projectID int64) (tasks []entity.Task, err error)
	UserTasks(ctx context.Context, userID int64) (tasks []entity.Task, err error)
}

type ProjectRepository interface {
	CreateProject(ctx context.Context, project entity.Project) (entity.Project, error)
	UserProjects(ctx context.Context, userID int64) (projects []entity.Project, err error)
	ProjectByID(ctx context.Context, id int64) (p entity.Project, err error)
	DeleteProject(ctx context.Context, projectID int64) error
	AddProjectMember(ctx context.Context, projectID int64, userID int64) error
}

type ProjectService struct {
	project ProjectRepository
	task    TaskRepository
}

func NewProjectRepository(project ProjectRepository, task TaskRepository) *ProjectService {
	return &ProjectService{
		project: project,
		task:    task,
	}
}

func (us *ProjectService) CreateProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	user := entity.AuthUser(ctx)

	project.UserID = user.ID
	project.CreatedAt = time.Now()

	project, err := us.project.CreateProject(ctx, project)
	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (us *ProjectService) ProjectByID(ctx context.Context, id int64) (entity.Project, error) {
	user := entity.AuthUser(ctx)

	project, err := us.project.ProjectByID(ctx, id)
	if err != nil {
		return entity.Project{}, err
	}

	if user.ID != project.UserID {
		return entity.Project{}, fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	return project, nil
}

func (us *ProjectService) UserProjects(ctx context.Context) ([]entity.Project, error) {
	user := entity.AuthUser(ctx)
	return us.project.UserProjects(ctx, user.ID)
}

func (us *ProjectService) DeleteProject(ctx context.Context, projectID int64) error {
	user := entity.AuthUser(ctx)

	project, err := us.project.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	if user.ID != project.UserID {
		return fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	err = us.project.DeleteProject(ctx, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (us *ProjectService) CreateTask(ctx context.Context, cTask entity.TaskToCreate) (entity.Task, error) {
	project, err := us.project.ProjectByID(ctx, cTask.ProjectID)
	if err != nil {
		return entity.Task{}, err
	}

	user := entity.AuthUser(ctx)

	if project.UserID != user.ID {
		return entity.Task{}, fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	task := entity.Task{
		Name:      cTask.Name,
		UserID:    user.ID,
		ProjectID: cTask.ProjectID,
		CreatedAt: time.Now(),
	}

	return us.task.CreateTask(ctx, task)
}

func (us *ProjectService) TaskByID(ctx context.Context, id int64) (entity.Task, error) {
	user := entity.AuthUser(ctx)

	task, err := us.task.TaskByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

	if user.ID != task.UserID {
		return entity.Task{}, fmt.Errorf("%w: not your task", entity.ErrForbidden)
	}

	return task, nil
}

func (us *ProjectService) ProjectTasks(ctx context.Context, projectID int64) ([]entity.Task, error) {
	project, err := us.project.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user := entity.AuthUser(ctx)

	if project.UserID != user.ID {
		return nil, fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	tasks, err := us.task.ProjectTasks(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (us *ProjectService) UserTasks(ctx context.Context) ([]entity.Task, error) {
	user := entity.AuthUser(ctx)

	tasks, err := us.task.UserTasks(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (us *ProjectService) AddProjectMember(ctx context.Context, projectID int64, userID int64) error {
	requester := entity.AuthUser(ctx)

	project, err := us.project.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	if requester.ID != project.UserID {
		return fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	err = us.project.AddProjectMember(ctx, projectID, userID)
	if err != nil {
		return err
	}

	return nil
}
