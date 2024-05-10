package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"restAPI/entity"
	"restAPI/service"
	"strconv"
	"time"
)

type Handler struct {
	us *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{us: userService}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, err)
		return
	}

	user, err = h.us.RegisterUser(ctx, user)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")

	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	err = h.us.DeleteUser(ctx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")

	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	user, err := h.us.UserByID(ctx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, user)
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.us.Users(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, users)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, err)
		return
	}

	sessionID, err := h.us.Login(ctx, user.Email, user.Password)
	if err != nil {
		sendError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID.String(),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		MaxAge:   24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project entity.Project

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		sendError(w, err)
		return
	}

	project, err = h.us.CreateProject(ctx, project)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, project)
}

func (h *Handler) UserProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projects, err := h.us.UserProjects(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, projects)
}

func (h *Handler) ProjectByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")
	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	project, err := h.us.ProjectByID(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, project)
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")
	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	err = h.us.DeleteProject(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cTask := entity.TaskToCreate{}

	err := json.NewDecoder(r.Body).Decode(&cTask)
	if err != nil {
		sendError(w, err)
		return
	}

	task, err := h.us.CreateTask(ctx, cTask)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, task)
}

func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")
	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	task, err := h.us.TaskByID(ctx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, task)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	ctx := r.Context()

	err := h.us.Verify(ctx, code)
	if err != nil {
		sendError(w, err)
		return
	}

	fmt.Fprint(w, "Verification Completed")
}

func (h *Handler) ProjectTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	qID := r.PathValue("project_id")

	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	tasks, err := h.us.ProjectTasks(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, tasks)
}

func (h *Handler) UserTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.us.UserTasks(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, tasks)
}

type AddProjectUserRequest struct {
	ProjectID int64 `json:"project_id"`
	UserID    int64 `json:"user_id"`
}

func (h *Handler) AddProjectUser(w http.ResponseWriter, r *http.Request) {
	var request AddProjectUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendError(w, err)
		return
	}

	ctx := r.Context()

	err = h.us.AddProjectMember(ctx, request.ProjectID, request.UserID)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ProjectUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	qID := r.PathValue("project_id")

	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	users, err := h.us.ProjectUsers(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, users)
}
