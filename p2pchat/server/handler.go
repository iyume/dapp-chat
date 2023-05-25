package server

import (
	"log"
	"net/http"
	"strings"
)

func newAccessHandler(handler http.Handler, token string) http.Handler {
	return accessHandler{token: token, next: handler}
}

type accessHandler struct {
	token string
	next  http.Handler
}

func (h accessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, *")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")

	// CORS preflight handle
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
	if token == "" {
		log.Panicln("no token to start server")
	}
	wrapped := newAccessHandler(srv, token)
	return wrapped
}
