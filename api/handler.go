package api

import (
	"encoding/json"
	"errors"
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
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, err)
		return
	}

	user, err = h.us.RegisterUser(user)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	qID := r.PathValue("id")

	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	err = h.us.DeleteUser(id)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserByID(w http.ResponseWriter, r *http.Request) {
	qID := r.PathValue("id")

	id, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, errors.New("'id' must be an integer"))
		return
	}

	user, err := h.us.UserByID(id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, user)
}

func (h *Handler) Users(w http.ResponseWriter, _ *http.Request) {
	users, err := h.us.Users()
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, users)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, err)
		return
	}

	sessionID, err := h.us.Login(user.Email, user.Password)
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
	cookie, err := r.Cookie("session_id")
	if err != nil {
		sendError(w, err)
		return
	}

	user, err := h.us.UserBySessionID(cookie.Value)
	if err != nil {
		sendError(w, err)
		return
	}

	project := entity.Project{
		Name:       "",
		OwnerID:    user.ID,
		OwnerEmail: user.Email,
		OwnerName:  user.Name,
		CreatedAt:  time.Now(),
	}

	project, err = h.us.CreateProject(project)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, project)
}

func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	qID := r.PathValue("owner_id")
	userID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, err)
		return
	}

	projects, err := h.us.Projects(userID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendResponse(w, projects)
}

func (h *Handler) ProjectByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	qID := r.PathValue("project_id")
	projectID, err := strconv.ParseInt(qID, 10, 64)
	if err != nil {
		sendError(w, err)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		sendError(w, err)
		return
	}

	user, err := h.us.UserBySessionID(cookie.Value)
	if err != nil {
		sendError(w, err)
		return
	}

	err = h.us.DeleteProject(user.ID, projectID)
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
