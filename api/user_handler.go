package api

import (
	"context"
	"errors"
	"net/http"
	"restAPI/entity"
	"strconv"
)

type UserService interface {
	UserByID(ctx context.Context, id int64) (entity.User, error)
	Users(ctx context.Context) ([]entity.User, error)
	ProjectUsers(ctx context.Context, projectID int64) ([]entity.User, error)

	DeleteUser(ctx context.Context, id int64) error
}

type UserHandler struct {
	user UserService
}

func NewUserHandler(user UserService) *UserHandler {
	return &UserHandler{user: user}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")

	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	err = h.user.DeleteUser(ctx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")

	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	user, err := h.user.UserByID(ctx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, user)
}

func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.user.Users(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, users)
}

func (h *UserHandler) ProjectUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	qID := r.PathValue("project_id")

	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	users, err := h.user.ProjectUsers(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, users)
}
