package api

import (
	"context"
	"log"
	"net/http"
	"restAPI/service"
)

type Middleware struct {
	us *service.UserService
}

func NewMiddleware(us *service.UserService) *Middleware {
	return &Middleware{us: us}
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

		user, err := mw.us.UserBySessionID(cookie.Value)
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
