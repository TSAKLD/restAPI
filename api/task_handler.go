package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"restAPI/entity"
	"strconv"
)

type TaskService interface {
	CreateTask(ctx context.Context, cTask entity.TaskToCreate) (entity.Task, error)
	TaskByID(ctx context.Context, id int64) (entity.Task, error)
	ProjectTasks(ctx context.Context, projectID int64) ([]entity.Task, error)
	UserTasks(ctx context.Context) ([]entity.Task, error)
}

type TaskHandler struct {
	task TaskService
}

func NewTaskHandler(task TaskService) *TaskHandler {
	return &TaskHandler{task: task}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cTask := entity.TaskToCreate{}

	err := json.NewDecoder(r.Body).Decode(&cTask)
	if err != nil {
		sendError(w, err)
		return
	}

	task, err := h.task.CreateTask(ctx, cTask)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, task)
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qID := r.PathValue("id")
	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	task, err := h.task.TaskByID(ctx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, task)
}

func (h *TaskHandler) ProjectTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	qID := r.PathValue("project_id")

	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	tasks, err := h.task.ProjectTasks(ctx, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, tasks)
}

func (h *TaskHandler) UserTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.task.UserTasks(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, tasks)
}
