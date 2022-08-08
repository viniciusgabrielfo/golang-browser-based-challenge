package handler

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/usecase/user"
)

func MakePublicHandlers(r *chi.Mux, userService user.Service, tokenAuth *jwtauth.JWTAuth) {
	r.Post("/register", handleRegisterUser(userService))
	r.Post("/auth", handleAuthUser(userService, tokenAuth))

	r.HandleFunc("/login", handleLoginPage())
	r.HandleFunc("/signup", handleSignupPage())
}

func handleLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("./chat/templates/login.html"))
		template.Execute(w, nil)
	}
}

func handleSignupPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("./chat/templates/signup.html"))
		template.Execute(w, nil)
	}
}

func handleRegisterUser(service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputRegister struct {
			Nick     string `json:"nick"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&inputRegister); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error decoding json user"))
			return
		}

		if err := service.CreateUser(inputRegister.Nick, inputRegister.Password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error adding user"))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func handleAuthUser(service user.Service, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputRegister struct {
			Nick     string `json:"nick"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&inputRegister); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error decoding json user"))
			return
		}

		userID, err := service.Auth(inputRegister.Nick, inputRegister.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error auth user, nick and password doesn't match"))
			return
		}

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": userID})

		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: tokenString,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenString))
	}
}
