package repository

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"restAPI/bootstrap"
	"restAPI/entity"
	"testing"
	"time"
)

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
		Name:      uuid.NewString(),
		Password:  uuid.NewString(),
		Email:     uuid.NewString(),
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}
	// Create user
	user, err = repo.CreateUser(user)
	require.NoError(t, err)

	// Get user by email && password
	user2, err := repo.UserByEmailAndPassword(user.Email, user.Password)
	require.NoError(t, err)

	user.Password = ""

	require.Equal(t, user, user2)

	// Get user by ID
	user2, err = repo.UserByID(user.ID)
	require.NoError(t, err)
	require.Equal(t, user, user2)

	// Get user by email
	user2, err = repo.UserByEmail(user.Email)
	require.NoError(t, err)
	require.Equal(t, user, user2)

	// Get users
	users, err := repo.Users()
	require.NoError(t, err)
	require.Contains(t, users, user)

	// Delete user
	err = repo.DeleteUser(user.ID)
	require.NoError(t, err)

	_, err = repo.UserByID(user.ID)
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

	_, err = repo.UserByEmailAndPassword(uuid.NewString(), uuid.NewString())
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = repo.UserByID(time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = repo.UserByEmail(uuid.NewString())
	require.ErrorIs(t, err, entity.ErrNotFound)

	db.Close()

	_, err = repo.CreateUser(entity.User{})
	require.Error(t, err)

	err = repo.DeleteUser(time.Now().UnixNano())
	require.Error(t, err)

	_, err = repo.Users()
	require.Error(t, err)

	_, err = repo.UserByEmailAndPassword(uuid.NewString(), uuid.NewString())
	require.Error(t, err)

	_, err = repo.UserByID(time.Now().UnixNano())
	require.Error(t, err)

	_, err = repo.UserByEmail(uuid.NewString())
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

	user, err = repo.CreateUser(user)
	require.NoError(t, err)

	actualProject := entity.Project{
		Name:      uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	// Create project
	actualProject, err = repo.CreateProject(actualProject)
	require.NoError(t, err)

	// User projects
	projects, err := repo.UserProjects(user.ID)
	require.NoError(t, err)
	require.Contains(t, projects, actualProject)

	// Project by ID
	expectedProject, err := repo.ProjectByID(actualProject.ID)
	require.NoError(t, err)
	require.Equal(t, expectedProject, actualProject)

	// Delete project
	err = repo.DeleteProject(actualProject.ID)
	require.NoError(t, err)

	_, err = repo.ProjectByID(actualProject.ID)
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

	_, err = repo.CreateProject(entity.Project{})
	require.Error(t, err)

	_, err = repo.ProjectByID(time.Now().UnixNano())
	require.ErrorIs(t, err, entity.ErrNotFound)

	db.Close()

	_, err = repo.UserProjects(time.Now().UnixNano())
	require.Error(t, err)

	err = repo.DeleteProject(time.Now().UnixNano())
	require.Error(t, err)
}

//func TestRepository_CreateTask(t *testing.T) {
//	cfg := &bootstrap.Config{
//		DBHost:     "localhost",
//		DBPort:     "5433",
//		DBUser:     "postgres",
//		DBPassword: "postgres",
//		DBName:     "postgres",
//	}
//
//	db, err := bootstrap.DBConnect(cfg)
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := New(db)
//
//	task := entity.Task{
//		ID:        time.Now().UnixNano(),
//		Name:      uuid.NewString(),
//		UserID:    time.Now().UnixNano(),
//		ProjectID: time.Now().UnixNano(),
//		CreatedAt: time.Time{},
//	}
//
//	task, err = repo.CreateTask(task)
//	require.NoError(t, err)
//}
