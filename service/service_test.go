package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"restAPI/entity"
	"testing"
	"time"
)

var testUser = entity.User{
	ID:         1,
	Name:       "dev",
	Password:   "dev",
	Email:      "dev",
	CreatedAt:  time.Now(),
	IsVerified: true,
}

type RepositoryMock struct {
	userByID func(ctx context.Context, id int64) (u entity.User, err error)
}

func (r RepositoryMock) CreateUser(ctx context.Context, u entity.User) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) DeleteUser(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) UserByID(ctx context.Context, id int64) (u entity.User, err error) {
	return r.userByID(ctx, id)
}

func (r RepositoryMock) UserByEmail(ctx context.Context, email string) (u entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) Users(ctx context.Context) (users []entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) UserByEmailAndPassword(ctx context.Context, email string, password string) (u entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) ProjectUsers(ctx context.Context, projectID int64) (users []entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) CreateSession(ctx context.Context, sessionID uuid.UUID, userID int64, createdAt time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) UserBySessionID(ctx context.Context, sessionID string) (u entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) SaveVerificationCode(ctx context.Context, code string, userID int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) VerifyUser(ctx context.Context, code string) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) CreateProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) UserProjects(ctx context.Context, userID int64) (projects []entity.Project, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) ProjectByID(ctx context.Context, id int64) (p entity.Project, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) DeleteProject(ctx context.Context, projectID int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) AddProjectMember(ctx context.Context, projectID int64, userID int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) CreateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) TaskByID(ctx context.Context, id int64) (t entity.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) ProjectTasks(ctx context.Context, projectID int64) (tasks []entity.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) UserTasks(ctx context.Context, userID int64) (tasks []entity.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func TestService_UserByID(t *testing.T) {
	mock := &RepositoryMock{}

	mock.userByID = func(ctx context.Context, id int64) (u entity.User, err error) {
		return testUser, nil
	}

	s := New(mock)

	user, err := s.UserByID(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, testUser, user)

	mock.userByID = func(ctx context.Context, id int64) (u entity.User, err error) {
		return entity.User{}, entity.ErrNotFound
	}

	user, err = s.UserByID(context.Background(), 2)
	require.ErrorIs(t, err, entity.ErrNotFound)

}
