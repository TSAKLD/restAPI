package service

import (
	"context"
	"fmt"
	"restAPI/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UserByID(ctx context.Context, id int64) (u entity.User, err error)
	UserByEmail(ctx context.Context, email string) (u entity.User, err error)
	Users(ctx context.Context) (users []entity.User, err error)
	ProjectUsers(ctx context.Context, projectID int64) (users []entity.User, err error)
}

type UserService struct {
	user    UserRepository
	auth    AuthRepository
	project ProjectRepository
}

func NewUserService(user UserRepository, auth AuthRepository, project ProjectRepository) *UserService {
	return &UserService{
		user:    user,
		auth:    auth,
		project: project,
	}
}

func (us *UserService) UserByID(ctx context.Context, id int64) (entity.User, error) {
	user, err := us.user.UserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id int64) error {
	_, err := us.user.UserByID(ctx, id)
	if err != nil {
		return err
	}

	err = us.user.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Users(ctx context.Context) ([]entity.User, error) {
	users, err := us.user.Users(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) ProjectUsers(ctx context.Context, projectID int64) ([]entity.User, error) {
	user := entity.AuthUser(ctx)

	projects, err := us.project.UserProjects(ctx, user.ID)
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

	users, err := us.user.ProjectUsers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
