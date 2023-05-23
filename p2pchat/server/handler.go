package server

import (
	"log"
	"net/http"
	"strings"
)

func newTokenHandler(handler http.Handler, token string) http.Handler {
	return tokenHandler{token: token, next: handler}
}

type tokenHandler struct {
	token string
	next  http.Handler
}

func (h tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "missing Authorization header", http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	} else {
		http.Error(w, "should be bearer token", http.StatusUnauthorized)
		return
	}
	if token != h.token {
		http.Error(w, "incorrect token", http.StatusUnauthorized)
		return
	}
	h.next.ServeHTTP(w, r)
}

func NewHTTPStack(srv http.Handler, token string) http.Handler {
	// TODO: CORS handler
	if token == "" {
		log.Panicln("no token to start server")
	}
	wrapped := newTokenHandler(srv, token)
	return wrapped
}
