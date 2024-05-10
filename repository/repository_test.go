package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"restAPI/bootstrap"
	"restAPI/entity"
	"testing"
	"time"
)

var eCtx = context.Background()

func TestRepository_CreateUser(t *testing.T) {
	cfg := &bootstrap.Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "postgres",
	}

	db, err := bootstrap.DBConnect(cfg)
	require.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)
	authRepo := NewAuthRepository(db)

	user := entity.User{
		Name:       uuid.NewString(),
		Password:   uuid.NewString(),
		Email:      uuid.NewString(),
		CreatedAt:  time.Now().UTC().Round(time.Millisecond),
		IsVerified: true,
	}
	// Create user
	user, err = userRepo.CreateUser(eCtx, user)
	require.NoError(t, err)

	// Get user by email && password
	user2, err := authRepo.UserByEmailAndPassword(eCtx, user.Email, user.Password)
	require.NoError(t, err)

	user.Password = ""

	require.Equal(t, user, user2)

	// Get user by ID
	user2, err = userRepo.UserByID(eCtx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user, user2)

	// Get user by email
	user2, err = userRepo.UserByEmail(eCtx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user, user2)

	// Get users
	users, err := userRepo.Users(eCtx)
	require.NoError(t, err)
	require.Contains(t, users, user)

	// Delete user
	err = userRepo.DeleteUser(eCtx, user.ID)
	require.NoError(t, err)

	_, err = userRepo.UserByID(eCtx, user.ID)
	require.ErrorIs(t, err, entity.ErrNotFound)
}

func TestRepository_Users_Error(t *testing.T) {
	cfg := &bootstrap.Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "postgres",
	}

	db, err := bootstrap.DBConnect(cfg)
	require.NoError(t, err)
	defer db.Close()

	authRepo := NewAuthRepository(db)
	userRepo := NewUserRepository(db)

	_, err = authRepo.UserByEmailAndPassword(eCtx, uuid.NewString(), uuid.NewString())
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = userRepo.UserByID(eCtx, time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = userRepo.UserByEmail(eCtx, uuid.NewString())
	require.ErrorIs(t, err, entity.ErrNotFound)

	db.Close()

	_, err = userRepo.CreateUser(eCtx, entity.User{})
	require.Error(t, err)

	err = userRepo.DeleteUser(eCtx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = userRepo.Users(eCtx)
	require.Error(t, err)

	_, err = authRepo.UserByEmailAndPassword(eCtx, uuid.NewString(), uuid.NewString())
	require.Error(t, err)

	_, err = userRepo.UserByID(eCtx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = userRepo.UserByEmail(eCtx, uuid.NewString())
	require.Error(t, err)
}

func TestRepository_CreateProject(t *testing.T) {
	cfg := &bootstrap.Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "postgres",
	}

	db, err := bootstrap.DBConnect(cfg)
	require.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)
	repo := NewProjectRepository(db)

	// Create user
	user := entity.User{
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	user, err = userRepo.CreateUser(eCtx, user)
	require.NoError(t, err)

	actualProject := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	// Create project
	actualProject, err = repo.CreateProject(eCtx, actualProject)
	require.NoError(t, err)

	// User projects
	projects, err := repo.UserProjects(eCtx, user.ID)
	require.NoError(t, err)
	require.Contains(t, projects, actualProject)

	// Project by ID
	expectedProject, err := repo.ProjectByID(eCtx, actualProject.ID)
	require.NoError(t, err)
	require.Equal(t, expectedProject, actualProject)

	// Delete project
	err = repo.DeleteProject(eCtx, actualProject.ID)
	require.NoError(t, err)

	_, err = repo.ProjectByID(eCtx, actualProject.ID)
	require.ErrorIs(t, err, entity.ErrNotFound)
}

func TestRepository_Projects_Error(t *testing.T) {
	cfg := &bootstrap.Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "postgres",
	}

	db, err := bootstrap.DBConnect(cfg)
	require.NoError(t, err)
	defer db.Close()

	repo := NewProjectRepository(db)

	_, err = repo.CreateProject(eCtx, entity.Project{})
	require.Error(t, err)

	_, err = repo.ProjectByID(eCtx, time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	db.Close()

	_, err = repo.UserProjects(eCtx, time.Now().UnixNano())
	require.Error(t, err)

	err = repo.DeleteProject(eCtx, time.Now().UnixNano())
	require.Error(t, err)
}

func TestRepository_CreateTask(t *testing.T) {
	cfg := &bootstrap.Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "postgres",
	}

	db, err := bootstrap.DBConnect(cfg)
	require.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)
	repo := NewProjectRepository(db)
	task := NewTaskRepository(db)

	// Create user
	user := entity.User{
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	user, err = userRepo.CreateUser(eCtx, user)
	require.NoError(t, err)

	actualProject := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualProject, err = repo.CreateProject(eCtx, actualProject)
	require.NoError(t, err)

	actualTask := entity.Task{
		Name:        uuid.NewString(),
		UserID:      user.ID,
		Description: uuid.NewString(),
		ProjectID:   actualProject.ID,
		CreatedAt:   time.Now().UTC().Round(time.Millisecond),
	}

	actualTask, err = task.CreateTask(eCtx, actualTask)
	require.NoError(t, err)

	expectedTask, err := task.TaskByID(eCtx, actualTask.ID)
	require.NoError(t, err)
	require.Equal(t, expectedTask, actualTask)

	_, err = task.TaskByID(eCtx, time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	////////////////////////

	user2 := entity.User{
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	user2, err = userRepo.CreateUser(eCtx, user2)
	require.NoError(t, err)

	actualProject2 := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user2.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualProject2, err = repo.CreateProject(eCtx, actualProject2)
	require.NoError(t, err)

	actualTask2 := entity.Task{
		Name:      uuid.NewString(),
		UserID:    user2.ID,
		ProjectID: actualProject2.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualTask2, err = task.CreateTask(eCtx, actualTask2)
	require.NoError(t, err)

	actualTasks, err := task.ProjectTasks(eCtx, actualProject.ID)
	require.NoError(t, err)
	require.Contains(t, actualTasks, actualTask)
	require.NotContains(t, actualTasks, actualTask2)

	actualTasks, err = task.UserTasks(eCtx, user.ID)
	require.NoError(t, err)
	require.Contains(t, actualTasks, actualTask)
	require.NotContains(t, actualTasks, actualTask2)

	db.Close()

	_, err = task.CreateTask(eCtx, entity.Task{})
	require.Error(t, err)

	_, err = task.TaskByID(eCtx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = task.ProjectTasks(eCtx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = task.UserTasks(eCtx, time.Now().UnixNano())
	require.Error(t, err)
}
