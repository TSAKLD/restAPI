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

type AuthRepository interface {
	UserByEmailAndPassword(ctx context.Context, email string, password string) (u entity.User, err error)
	CreateSession(ctx context.Context, sessionID uuid.UUID, userID int64, createdAt time.Time) error
	UserBySessionID(ctx context.Context, sessionID string) (u entity.User, err error)
	SaveVerificationCode(ctx context.Context, code string, userID int64) error
	VerifyUser(ctx context.Context, code string) error
}

type AuthService struct {
	auth AuthRepository
	user UserRepository
}

func NewAuthService(auth AuthRepository) *AuthService {
	return &AuthService{auth: auth}
}

func (us *AuthService) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
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

func (us *AuthService) Login(ctx context.Context, email string, password string) (uuid.UUID, error) {
	user, err := us.auth.UserByEmailAndPassword(ctx, email, password)
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

	err = us.auth.CreateSession(ctx, sessionID, user.ID, createdAt)
	if err != nil {
		return uuid.UUID{}, err
	}

	return sessionID, nil
}

func (us *AuthService) UserBySessionID(ctx context.Context, sessionID string) (entity.User, error) {
	return us.auth.UserBySessionID(ctx, sessionID)
}

func (us *AuthService) Verify(ctx context.Context, code string) error {
	return us.auth.VerifyUser(ctx, code)
}

func (us *AuthService) SendVerificationLink(ctx context.Context, code string, email string) error {
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
