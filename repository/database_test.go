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

var Emptyctx = context.Background()

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

	repo := New(db)

	user := entity.User{
		Name:       uuid.NewString(),
		Password:   uuid.NewString(),
		Email:      uuid.NewString(),
		CreatedAt:  time.Now().UTC().Round(time.Millisecond),
		IsVerified: true,
	}
	// Create user
	user, err = repo.CreateUser(Emptyctx, user)
	require.NoError(t, err)

	// Get user by email && password
	user2, err := repo.UserByEmailAndPassword(Emptyctx, user.Email, user.Password)
	require.NoError(t, err)

	user.Password = ""

	require.Equal(t, user, user2)

	// Get user by ID
	user2, err = repo.UserByID(Emptyctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user, user2)

	// Get user by email
	user2, err = repo.UserByEmail(Emptyctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user, user2)

	// Get users
	users, err := repo.Users(Emptyctx)
	require.NoError(t, err)
	require.Contains(t, users, user)

	// Delete user
	err = repo.DeleteUser(Emptyctx, user.ID)
	require.NoError(t, err)

	_, err = repo.UserByID(Emptyctx, user.ID)
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

	repo := New(db)

	_, err = repo.UserByEmailAndPassword(Emptyctx, uuid.NewString(), uuid.NewString())
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = repo.UserByID(Emptyctx, time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = repo.UserByEmail(Emptyctx, uuid.NewString())
	require.ErrorIs(t, err, entity.ErrNotFound)

	db.Close()

	_, err = repo.CreateUser(Emptyctx, entity.User{})
	require.Error(t, err)

	err = repo.DeleteUser(Emptyctx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = repo.Users(Emptyctx)
	require.Error(t, err)

	_, err = repo.UserByEmailAndPassword(Emptyctx, uuid.NewString(), uuid.NewString())
	require.Error(t, err)

	_, err = repo.UserByID(Emptyctx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = repo.UserByEmail(Emptyctx, uuid.NewString())
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

	repo := New(db)

	// Create user
	user := entity.User{
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	user, err = repo.CreateUser(Emptyctx, user)
	require.NoError(t, err)

	actualProject := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	// Create project
	actualProject, err = repo.CreateProject(Emptyctx, actualProject)
	require.NoError(t, err)

	// User projects
	projects, err := repo.UserProjects(Emptyctx, user.ID)
	require.NoError(t, err)
	require.Contains(t, projects, actualProject)

	// Project by ID
	expectedProject, err := repo.ProjectByID(Emptyctx, actualProject.ID)
	require.NoError(t, err)
	require.Equal(t, expectedProject, actualProject)

	// Delete project
	err = repo.DeleteProject(Emptyctx, actualProject.ID)
	require.NoError(t, err)

	_, err = repo.ProjectByID(Emptyctx, actualProject.ID)
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

	repo := New(db)

	_, err = repo.CreateProject(Emptyctx, entity.Project{})
	require.Error(t, err)

	_, err = repo.ProjectByID(Emptyctx, time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	db.Close()

	_, err = repo.UserProjects(Emptyctx, time.Now().UnixNano())
	require.Error(t, err)

	err = repo.DeleteProject(Emptyctx, time.Now().UnixNano())
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

	repo := New(db)

	// Create user
	user := entity.User{
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	user, err = repo.CreateUser(Emptyctx, user)
	require.NoError(t, err)

	actualProject := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualProject, err = repo.CreateProject(Emptyctx, actualProject)
	require.NoError(t, err)

	actualTask := entity.Task{
		Name:      uuid.NewString(),
		UserID:    user.ID,
		ProjectID: actualProject.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualTask, err = repo.CreateTask(Emptyctx, actualTask)
	require.NoError(t, err)

	expectedTask, err := repo.TaskByID(Emptyctx, actualTask.ID)
	require.NoError(t, err)
	require.Equal(t, expectedTask, actualTask)

	_, err = repo.TaskByID(Emptyctx, time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	////////////////////////

	user2 := entity.User{
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	user2, err = repo.CreateUser(Emptyctx, user2)
	require.NoError(t, err)

	actualProject2 := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user2.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualProject2, err = repo.CreateProject(Emptyctx, actualProject2)
	require.NoError(t, err)

	actualTask2 := entity.Task{
		Name:      uuid.NewString(),
		UserID:    user2.ID,
		ProjectID: actualProject2.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	actualTask2, err = repo.CreateTask(Emptyctx, actualTask2)
	require.NoError(t, err)

	actualTasks, err := repo.ProjectTasks(Emptyctx, actualProject.ID)
	require.NoError(t, err)
	require.Contains(t, actualTasks, actualTask)
	require.NotContains(t, actualTasks, actualTask2)

	actualTasks, err = repo.UserTasks(Emptyctx, user.ID)
	require.NoError(t, err)
	require.Contains(t, actualTasks, actualTask)
	require.NotContains(t, actualTasks, actualTask2)

	db.Close()

	_, err = repo.CreateTask(Emptyctx, entity.Task{})
	require.Error(t, err)

	_, err = repo.TaskByID(Emptyctx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = repo.ProjectTasks(Emptyctx, time.Now().UnixNano())
	require.Error(t, err)

	_, err = repo.UserTasks(Emptyctx, time.Now().UnixNano())
	require.Error(t, err)
}
