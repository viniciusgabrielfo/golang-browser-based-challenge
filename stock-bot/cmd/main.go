package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/viniciusgabrielfo/golang-browser-based-challenge/stock-bot/internal"
	"go.uber.org/zap"
)

var (
	chatAddr   string
	chatWsPath string
	log        *zap.SugaredLogger
)

func setupFlags() {
	flag.StringVar(&chatAddr, "chat-addr", ":8000", "http chat ws service address")
	flag.StringVar(&chatWsPath, "chat-ws", "/ws", "http chat ws service path")

	flag.Parse()
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log = logger.Sugar()

	setupFlags()

	wsClient, err := internal.NewWebSocketClient("ws", chatAddr, chatWsPath, make(chan []byte), log)
	if err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()

	wsClient.Write("BOT ENTOUR")

	<-stop
}
