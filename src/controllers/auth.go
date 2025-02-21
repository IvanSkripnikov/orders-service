package controllers

import (
	"net/http"

	"authenticator/helpers"
)

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		helpers.Register(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/register")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		helpers.Login(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/login")
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		helpers.Auth(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/auth")
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.SignIn(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/signin")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		helpers.Logout(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/logout")
	}
}

func Sessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.Sessions(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/sessions")
	}
}
