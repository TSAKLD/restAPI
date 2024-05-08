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
		CreatedAt: time.Now(),
	}

	user, err = repo.CreateUser(user)
	require.NoError(t, err)

	user2, err := repo.UserByID(user.ID)
	require.NoError(t, err)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Email, user2.Email)
	require.NotEmpty(t, user2.CreatedAt)
}
