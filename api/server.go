package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	port    string
	router  *http.ServeMux
	taskHdr *TaskHandler
	projHdr *ProjectHandler
	userHdr *UserHandler
	authHdr *AuthHandler
	mw      *Middleware
}

// NewServer returns http router to work with.
func NewServer(t *TaskHandler, p *ProjectHandler, u *UserHandler, a *AuthHandler, port string, mw *Middleware) *Server {
	return &Server{
		port:    port,
		router:  http.NewServeMux(),
		taskHdr: t,
		projHdr: p,
		userHdr: u,
		authHdr: a,
		mw:      mw,
	}
}

// setRoutes activating handlers and sets routes for http router.
func (s *Server) setRoutes() {
	// user routes
	s.router.HandleFunc("DELETE /users/{id}", s.userHdr.DeleteUser)
	//s.router.HandleFunc("DELETE /users/{id}", s.h.EditUser)
	s.router.HandleFunc("GET /users/{id}", s.userHdr.UserByID)
	s.router.HandleFunc("GET /users", s.userHdr.Users)
	s.router.Handle("GET /projects/{project_id}/users", s.mw.Auth(s.userHdr.ProjectUsers))

	// auth routes
	s.router.HandleFunc("POST /users", s.authHdr.CreateUser)
	s.router.HandleFunc("GET /users/verify", s.authHdr.Verify)
	s.router.HandleFunc("POST /signin", s.authHdr.SignIn)

	// project routes
	s.router.Handle("POST /projects", s.mw.Auth(s.projHdr.CreateProject))
	s.router.Handle("DELETE /projects/{id}", s.mw.Auth(s.projHdr.DeleteProject))
	s.router.Handle("GET /projects", s.mw.Auth(s.projHdr.UserProjects))
	s.router.Handle("GET /projects/{id}", s.mw.Auth(s.projHdr.ProjectByID))
	//s.router.HandleFunc("POST /projects", s.h.EditProject)
	s.router.Handle("POST /projects/users", s.mw.Auth(s.projHdr.AddProjectUser))

	// task routes
	s.router.Handle("POST /tasks", s.mw.Auth(s.taskHdr.CreateTask))
	s.router.Handle("GET /tasks/{id}", s.mw.Auth(s.taskHdr.TaskByID))
	s.router.Handle("GET /projects/{project_id}/tasks", s.mw.Auth(s.taskHdr.ProjectTasks))
	s.router.Handle("GET /tasks", s.mw.Auth(s.taskHdr.UserTasks))
}

func (s *Server) Start() error {
	s.setRoutes()

	fmt.Println("Server is listening... at post:", s.port)

	return http.ListenAndServe(":"+s.port, s.mw.Log(s.router))
}
