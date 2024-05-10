package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
