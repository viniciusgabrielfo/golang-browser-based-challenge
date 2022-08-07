package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/viniciusgabrielfo/golang-browser-based-challenge/stock-bot/internal"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/stock-bot/pkg/rabbitmq"
	"go.uber.org/zap"
)

const ExchangeChat = "chat"

var (
	rabbitAddr       string
	rabbitUser       string
	rabbitPass       string
	rabbitRoutingKey string
	chatAddr         string
	chatWsPath       string
	chatAuthNick     string
	chatAuthPass     string
	log              *zap.SugaredLogger
)

func setupFlags() {
	flag.StringVar(&chatAddr, "chat-addr", "localhost:8000", "http chat ws service address")
	flag.StringVar(&chatWsPath, "chat-ws", "/ws", "http chat ws service path")
	flag.StringVar(&rabbitAddr, "rabbit_addr", "localhost:5672", "rabbitmq host addr")
	flag.StringVar(&rabbitUser, "rabbit_user", "guest", "rabbitmq user")
	flag.StringVar(&rabbitPass, "rabbit_password", "guest", "rabbitmq password")
	flag.StringVar(&rabbitRoutingKey, "rabbit_key", "chat", "rabbitmq binding routing key")
	flag.StringVar(&chatAuthNick, "chat_nick", "stock-bot", "chat nick used to authentication in a chatroom")
	flag.StringVar(&chatAuthPass, "chat_pass", "stockbotpass123", "chat password used to authentication in a chatroom")

	flag.Parse()
}

func init() {
	setupFlags()
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	log = logger.Sugar()

	stockBot := internal.NewStockBot(chatAuthNick, configStockBotProducer())
	go stockBot.Start()
	log.Info("stock bot commands listener started...")

	wsClient, err := internal.NewWebSocketClient("ws", chatAddr, chatWsPath, chatAuthNick, chatAuthPass, stockBot.GetListener(), log)
	if err != nil {
		log.Fatal(err)
	}

	<-stop
	log.Info("starting graceful shutdown...")

	stockBot.Stop()
	wsClient.Close()

	log.Info("stock-bot finished")
}

func configStockBotProducer() rabbitmq.Producer {
	amqpConn, err := rabbitmq.NewAmqpConnManager(rabbitmq.AmqpConfig{
		User:     rabbitUser,
		Password: rabbitPass,
		Host:     rabbitAddr,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := amqpConn.ExchangeDeclare(ExchangeChat, "direct"); err != nil {
		log.Fatal(err)
	}

	return amqpConn.CreateProducer(ExchangeChat, rabbitRoutingKey)
}
