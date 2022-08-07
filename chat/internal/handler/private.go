package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/usecase/user"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/usecase/wsclient"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func MakePrivateHandlers(r *chi.Mux, chatroom *entity.Chatroom, userService user.Service, tokenAuth *jwtauth.JWTAuth, logger *zap.SugaredLogger) {
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			user, err := userService.GetUser(uuid.MustParse(claims["user_id"].(string)))
			if err != nil {
				if !errors.Is(err, entity.ErrNotFoundEntity) {
					logger.Error(err)
				}
				return
			}

			w.Write([]byte(fmt.Sprintf("protected area. hi %v", user.Nick)))
		})

		r.HandleFunc("/ws", handleWebSocketConn(chatroom, userService, logger))
	})
}

func handleWebSocketConn(chatroom *entity.Chatroom, userService user.Service, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(err)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		user, err := userService.GetUser(uuid.MustParse(claims["user_id"].(string)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error identity user logged"))
			return
		}

		clientService := wsclient.NewService(entity.NewClient(conn, user, chatroom), logger)

		chatroom.GetRegisterChan() <- clientService.GetClient()
		defer func() {
			chatroom.GetUnregisterChan() <- clientService.GetClient()
		}()

		go clientService.Write()
		clientService.Read()
	}
}
