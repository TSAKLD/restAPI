package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"restAPI/entity"
	"strconv"
)

type ProjectService interface {
	CreateProject(ctx context.Context, project entity.Project) (entity.Project, error)
	ProjectByID(ctx context.Context, id int64) (entity.Project, error)
	UserProjects(ctx context.Context) ([]entity.Project, error)
	DeleteProject(ctx context.Context, projectID int64) error
	AddProjectMember(ctx context.Context, projectID int64, userID int64) error
}

type ProjectHandler struct {
	project ProjectService
}

func NewProjectHandler(project ProjectService) *ProjectHandler {
	return &ProjectHandler{project: project}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project entity.Project

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		sendError(w, err)
		return
	}

	project, err = h.project.CreateProject(ctx, project)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, project)
}

func (h *ProjectHandler) UserProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projects, err := h.project.UserProjects(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, projects)
}

func (h *ProjectHandler) ProjectByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")
	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	project, err := h.project.ProjectByID(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")
	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	err = h.project.DeleteProject(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type AddProjectUserRequest struct {
	ProjectID int64 `json:"project_id"`
	UserID    int64 `json:"user_id"`
}

func (h *ProjectHandler) AddProjectUser(w http.ResponseWriter, r *http.Request) {
	var request AddProjectUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendError(w, err)
		return
	}

	ctx := r.Context()

	err = h.project.AddProjectMember(ctx, request.ProjectID, request.UserID)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
