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
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, u entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UserByID(ctx context.Context, id int64) (u entity.User, err error)
	UserByEmail(ctx context.Context, email string) (u entity.User, err error)
	Users(ctx context.Context) (users []entity.User, err error)
	UserByEmailAndPassword(ctx context.Context, email string, password string) (u entity.User, err error)
	ProjectUsers(ctx context.Context, projectID int64) (users []entity.User, err error)
	CreateSession(ctx context.Context, sessionID uuid.UUID, userID int64, createdAt time.Time) error
	UserBySessionID(ctx context.Context, sessionID string) (u entity.User, err error)
	SaveVerificationCode(ctx context.Context, code string, userID int64) error
	VerifyUser(ctx context.Context, code string) error
	CreateProject(ctx context.Context, project entity.Project) (entity.Project, error)
	UserProjects(ctx context.Context, userID int64) (projects []entity.Project, err error)
	ProjectByID(ctx context.Context, id int64) (p entity.Project, err error)
	DeleteProject(ctx context.Context, projectID int64) error
	AddProjectMember(ctx context.Context, projectID int64, userID int64) error
	CreateTask(ctx context.Context, t entity.Task) (entity.Task, error)
	TaskByID(ctx context.Context, id int64) (t entity.Task, err error)
	ProjectTasks(ctx context.Context, projectID int64) (tasks []entity.Task, err error)
	UserTasks(ctx context.Context, userID int64) (tasks []entity.Task, err error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// user manipulations

func (us *Service) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
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

func (us *Service) UserByID(ctx context.Context, id int64) (entity.User, error) {
	user, err := us.repo.UserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (us *Service) DeleteUser(ctx context.Context, id int64) error {
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

func (us *Service) Users(ctx context.Context) ([]entity.User, error) {
	users, err := us.repo.Users(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *Service) Login(ctx context.Context, email string, password string) (uuid.UUID, error) {
	user, err := us.repo.UserByEmailAndPassword(ctx, email, password)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return uuid.UUID{}, entity.ErrUnauthorized
		}

		return uuid.UUID{}, err
	}

	if !user.IsVerified {
		return uuid.UUID{}, fmt.Errorf("%w: not verified, check your email", entity.ErrUnauthorized)
	}

	sessionID := uuid.New()

	createdAt := time.Now()

	err = us.repo.CreateSession(ctx, sessionID, user.ID, createdAt)
	if err != nil {
		return uuid.UUID{}, err
	}

	return sessionID, nil
}

func (us *Service) UserBySessionID(ctx context.Context, sessionID string) (entity.User, error) {
	return us.repo.UserBySessionID(ctx, sessionID)
}

// project manipulations

func (us *Service) CreateProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	user := entity.AuthUser(ctx)

	project.UserID = user.ID
	project.CreatedAt = time.Now()

	project, err := us.repo.CreateProject(ctx, project)
	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (us *Service) ProjectByID(ctx context.Context, id int64) (entity.Project, error) {
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

func (us *Service) UserProjects(ctx context.Context) ([]entity.Project, error) {
	user := entity.AuthUser(ctx)
	return us.repo.UserProjects(ctx, user.ID)
}

func (us *Service) DeleteProject(ctx context.Context, projectID int64) error {
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

func (us *Service) CreateTask(ctx context.Context, cTask entity.TaskToCreate) (entity.Task, error) {
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

func (us *Service) TaskByID(ctx context.Context, id int64) (entity.Task, error) {
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

func (us *Service) SendVerificationLink(ctx context.Context, code string, email string) error {
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

func (us *Service) Verify(ctx context.Context, code string) error {
	return us.repo.VerifyUser(ctx, code)
}

func (us *Service) ProjectTasks(ctx context.Context, projectID int64) ([]entity.Task, error) {
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

func (us *Service) UserTasks(ctx context.Context) ([]entity.Task, error) {
	user := entity.AuthUser(ctx)

	tasks, err := us.repo.UserTasks(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (us *Service) AddProjectMember(ctx context.Context, projectID int64, userID int64) error {
	requester := entity.AuthUser(ctx)

	project, err := us.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	if requester.ID != project.UserID {
		return fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	err = us.repo.AddProjectMember(ctx, projectID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (us *Service) ProjectUsers(ctx context.Context, projectID int64) ([]entity.User, error) {
	user := entity.AuthUser(ctx)

	projects, err := us.repo.UserProjects(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result := false

	for _, v := range projects {
		if projectID == v.ID {
			result = true
			break
		}
	}

	if !result {
		return nil, fmt.Errorf("%w: not your project", entity.ErrForbidden)
	}

	users, err := us.repo.ProjectUsers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
