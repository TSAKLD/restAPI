package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"restAPI/entity"
	"restAPI/repository"
	"time"
)

type UserService struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

// user manipulations

func (us *UserService) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	_, err := us.repo.UserByEmail(ctx, user.Email)
	if err == nil {
		return entity.User{}, fmt.Errorf("email %s already exist", user.Email)
	}

	user.CreatedAt = time.Now()

	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	user.Password = ""
	/////

	code := uuid.NewString()

	err = us.repo.SaveVerificationCode(ctx, code, user.ID)
	if err != nil {
		return entity.User{}, err
	}

	err = us.SendVerificationLink(ctx, code, user.Email)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (us *UserService) UserByID(ctx context.Context, id int64) (entity.User, error) {
	user, err := us.repo.UserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id int64) error {
	_, err := us.repo.UserByID(ctx, id)
	if err != nil {
		return err
	}

	err = us.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Users(ctx context.Context) ([]entity.User, error) {
	users, err := us.repo.Users(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) Login(ctx context.Context, email string, password string) (uuid.UUID, error) {
	user, err := us.repo.UserByEmailAndPassword(ctx, email, password)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return uuid.UUID{}, entity.ErrUnauthorized
		}

		return uuid.UUID{}, err
	}

	sessionID := uuid.New()

	createdAt := time.Now()

	err = us.repo.CreateSession(ctx, sessionID, user.ID, createdAt)
	if err != nil {
		return uuid.UUID{}, err
	}

	return sessionID, nil
}

func (us *UserService) UserBySessionID(ctx context.Context, sessionID string) (entity.User, error) {
	return us.repo.UserBySessionID(ctx, sessionID)
}

// project manipulations

func (us *UserService) CreateProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	user := entity.AuthUser(ctx)

	project.UserID = user.ID
	project.CreatedAt = time.Now()

	project, err := us.repo.CreateProject(ctx, project)
	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (us *UserService) ProjectByID(ctx context.Context, id int64) (entity.Project, error) {
	user := entity.AuthUser(ctx)

	project, err := us.repo.ProjectByID(ctx, id)
	if err != nil {
		return entity.Project{}, err
	}

	if user.ID != project.UserID {
		return entity.Project{}, fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	return project, nil
}

func (us *UserService) UserProjects(ctx context.Context) ([]entity.Project, error) {
	user := entity.AuthUser(ctx)
	return us.repo.UserProjects(ctx, user.ID)
}

func (us *UserService) DeleteProject(ctx context.Context, projectID int64) error {
	user := entity.AuthUser(ctx)

	project, err := us.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	if user.ID != project.UserID {
		return fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	err = us.repo.DeleteProject(ctx, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) CreateTask(ctx context.Context, cTask entity.TaskToCreate) (entity.Task, error) {
	project, err := us.repo.ProjectByID(ctx, cTask.ProjectID)
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

	return us.repo.CreateTask(ctx, task)
}

func (us *UserService) TaskByID(ctx context.Context, id int64) (entity.Task, error) {
	user := entity.AuthUser(ctx)

	task, err := us.repo.TaskByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

	if user.ID != task.UserID {
		return entity.Task{}, fmt.Errorf("%w: not your task", entity.ErrForbidden)
	}

	return task, nil
}

func (us *UserService) SendVerificationLink(ctx context.Context, code string, email string) error {
	message := map[string]interface{}{
		"subject":  "Verification",
		"receiver": email,
		"message":  fmt.Sprintf("Your Verification link is:http://localhost:8080/users/verify?code=%s", code),
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	response, err := http.Post("http://localhost:8090/mail", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (us *UserService) Verify(ctx context.Context, code string) error {
	return us.repo.Verify(ctx, code)
}

func (us *UserService) ProjectTasks(ctx context.Context, projectID int64) ([]entity.Task, error) {
	project, err := us.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	user := entity.AuthUser(ctx)

	if project.UserID != user.ID {
		return []entity.Task{}, fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	tasks, err := us.repo.ProjectTasks(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (us *UserService) UserTasks(ctx context.Context) ([]entity.Task, error) {
	user := entity.AuthUser(ctx)

	tasks, err := us.repo.UserTasks(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
