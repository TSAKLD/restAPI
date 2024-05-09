package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	port   string
	router *http.ServeMux
	h      *Handler
	mw     *Middleware
}

// NewServer returns http router to work with.
func NewServer(h *Handler, port string, mw *Middleware) *Server {
	return &Server{
		port:   port,
		router: http.NewServeMux(),
		h:      h,
		mw:     mw,
	}
}

// setRoutes activating handlers and sets routes for http router.
func (s *Server) setRoutes() {
	// user routes
	s.router.HandleFunc("POST /users", s.h.CreateUser)
	s.router.HandleFunc("DELETE /users/{id}", s.h.DeleteUser)
	//s.router.HandleFunc("DELETE /users/{id}", s.h.EditUser)
	s.router.HandleFunc("GET /users/{id}", s.h.UserByID)
	s.router.HandleFunc("GET /users", s.h.Users)

	s.router.HandleFunc("POST /signin", s.h.SignIn)

	// project routes
	s.router.Handle("POST /projects", s.mw.Auth(s.h.CreateProject))
	s.router.Handle("DELETE /projects/{id}", s.mw.Auth(s.h.DeleteProject))
	s.router.Handle("GET /projects", s.mw.Auth(s.h.UserProjects))
	s.router.Handle("GET /projects/{id}", s.mw.Auth(s.h.ProjectByID))
	//s.router.HandleFunc("POST /projects", s.h.EditProject)

	// task routes
	s.router.Handle("POST /tasks", s.mw.Auth(s.h.CreateTask))
	s.router.Handle("GET /tasks/{id}", s.mw.Auth(s.h.TaskByID))
}

func (s *Server) Start() error {
	s.setRoutes()

	fmt.Println("Server is listening... at post:", s.port)

	return http.ListenAndServe(":"+s.port, s.mw.Log(s.router))
}
