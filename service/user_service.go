package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"restAPI/entity"
	"time"
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
	user UserRepository
	auth AuthRepository
}

func NewUserService(user UserRepository, auth AuthRepository) *UserService {
	return &UserService{
		user: user,
		auth: auth,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	_, err := us.user.UserByEmail(ctx, user.Email)
	if err == nil {
		return entity.User{}, fmt.Errorf("email %s already exist", user.Email)
	}

	user.CreatedAt = time.Now()

	user, err = us.user.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	user.Password = ""
	/////

	code := uuid.NewString()

	err = us.auth.SaveVerificationCode(ctx, code, user.ID)
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

func (us *ProjectService) ProjectUsers(ctx context.Context, projectID int64) ([]entity.User, error) {
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
