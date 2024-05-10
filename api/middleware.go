package api

import (
	"context"
	"log"
	"net/http"
)

type Middleware struct {
	project ProjectService
	auth    AuthService
	user    UserService
}

func NewMiddleware(project ProjectService, auth AuthService, user UserService) *Middleware {
	return &Middleware{
		project: project,
		auth:    auth,
		user:    user,
	}
}

func (mw *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (mw *Middleware) Auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sendError(w, err)
			return
		}

		user, err := mw.auth.UserBySessionID(context.Background(), cookie.Value)
		if err != nil {
			sendError(w, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
