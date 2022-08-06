package main

import (
	"flag"
	"net/http"

	"github.com/viniciusgabrielfo/golang-browser-based-challenge/internal/entity"
	"go.uber.org/zap"
)

var (
	addr string
	log  *zap.SugaredLogger
)

func setupFlags() {
	flag.StringVar(&addr, "addr", ":8000", "http address")

	flag.Parse()
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log = logger.Sugar()

	setupFlags()

	chatroom := entity.NewChatroom(log)
	go chatroom.Start()

	http.HandleFunc("/ws", chatroom.HandleWebSocketConn)

	log.Info("starting HTTP server on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
