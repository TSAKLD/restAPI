package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	port   string
	router *http.ServeMux
	h      *Handler
}

// NewServer returns http router to work with.
func NewServer(h *Handler, port string) *Server {
	return &Server{
		port:   port,
		router: http.NewServeMux(),
		h:      h,
	}
}

// setRoutes activating handlers and sets routes for http router.
func (s *Server) setRoutes() {
	s.router.HandleFunc("POST /users", s.h.CreateUser)
	s.router.HandleFunc("DELETE /users/{id}", s.h.DeleteUser)
	s.router.HandleFunc("GET /users/{id}", s.h.UserByID)
	s.router.HandleFunc("GET /users", s.h.Users)

	s.router.HandleFunc("POST /signin", s.h.SignIn)
	s.router.HandleFunc("POST /projects", s.h.CreateProject)
}

func (s *Server) Start() error {
	s.setRoutes()

	fmt.Println("Server is listening... at post:", s.port)

	return http.ListenAndServe(":"+s.port, s.router)
}
