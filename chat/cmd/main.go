package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/handler"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/repository"
	chatroomsvc "github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/usecase/chatroom"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/usecase/user"
	"go.uber.org/zap"
)

const ExchangeChat = "chat"
const TopicChatOutbound = "chat.outbound"

var (
	addr             string
	rabbitAddr       string
	rabbitRoutingKey string
	jwtSecret        string
	stockBotPassword string
	log              *zap.SugaredLogger
)

func setupFlags() {
	flag.StringVar(&addr, "addr", ":8000", "http address")
	flag.StringVar(&rabbitAddr, "rabbit_addr", "amqp://guest:guest@localhost:5672/", "rabbitmq addr")
	flag.StringVar(&rabbitRoutingKey, "rabbit_key", "chat", "rabbitmq binding routing key")
	flag.StringVar(&jwtSecret, "jwt_secret", "defaultsecret", "jwt secret used to generate jwt tokens")
	flag.StringVar(&stockBotPassword, "stockbot_password", "stockbotpass123", "password used by stock-bot to auth in chatroom")

	flag.Parse()
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log = logger.Sugar()

	setupFlags()

	conn, err := amqp.Dial(rabbitAddr)
	if err != nil {
		log.Fatal(conn)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	if err := ch.ExchangeDeclare(
		ExchangeChat,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal(err)
	}

	queue, err := ch.QueueDeclare(
		TopicChatOutbound,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = ch.QueueBind(
		queue.Name,
		rabbitRoutingKey,
		ExchangeChat,
		false,
		nil,
	); err != nil {
		log.Fatal(err)
	}

	chatroom := entity.NewChatroom()
	userInMemory := repository.NewUserInMemory()

	chatroomService := chatroomsvc.NewService(chatroom, log)
	userService := user.NewService(userInMemory, log)

	createStockBotUser(userInMemory)

	go chatroomService.Start()

	consumer := internal.NewConsumer(ch, "consumer1", TopicChatOutbound)
	if err := consumer.Start(chatroom.GetBroadcastChan()); err != nil {
		log.Fatal(err)
	}

	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	r := chi.NewRouter()
	handler.MakePrivateHandlers(r, chatroom, userService, tokenAuth, log)
	handler.MakePublicHandlers(r, userService, tokenAuth)

	log.Info("starting HTTP server on " + addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	if err := consumer.Close(); err != nil {
		log.Fatal(err)
	}
}

func createStockBotUser(userRepo user.Repository) {
	if err := userRepo.Create(entity.NewUser("stock-bot", stockBotPassword)); err != nil {
		log.Error(err)
	}
}
