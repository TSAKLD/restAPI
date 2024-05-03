package service

import (
	"errors"
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

func (us *UserService) RegisterUser(user entity.User) (entity.User, error) {
	_, err := us.repo.UserByEmail(user.Email)
	if err == nil {
		return entity.User{}, errors.New("email %v already exist")
	}

	user.CreatedAt = time.Now()

	user, err = us.repo.CreateUser(user)
	if err != nil {
		return entity.User{}, err
	}

	user.Password = ""

	return user, nil
}

func (us *UserService) UserByID(id int64) (entity.User, error) {
	user, err := us.repo.UserByID(id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (us *UserService) DeleteUser(id int64) error {
	_, err := us.repo.UserByID(id)
	if err != nil {
		return err
	}

	err = us.repo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Users() ([]entity.User, error) {
	users, err := us.repo.Users()
	if err != nil {
		return nil, err
	}

	return users, nil
}
