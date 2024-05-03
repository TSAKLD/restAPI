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

	fmt.Println(user)
}
