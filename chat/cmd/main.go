package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/handler"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/pkg/rabbitmq"
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
	rabbitUser       string
	rabbitPass       string
	rabbitRoutingKey string
	jwtSecret        string
	stockBotPassword string
)

var (
	log       *zap.SugaredLogger
	tokenAuth *jwtauth.JWTAuth
)

func setupFlags() {
	flag.StringVar(&addr, "addr", ":8000", "http address")
	flag.StringVar(&rabbitAddr, "rabbit_addr", "localhost:5672", "rabbitmq host addr")
	flag.StringVar(&rabbitUser, "rabbit_user", "guest", "rabbitmq user")
	flag.StringVar(&rabbitPass, "rabbit_password", "guest", "rabbitmq password")
	flag.StringVar(&rabbitRoutingKey, "rabbit_key", "chat", "rabbitmq binding routing key")
	flag.StringVar(&jwtSecret, "jwt_secret", "defaultsecret", "jwt secret used to generate jwt tokens")
	flag.StringVar(&stockBotPassword, "stockbot_password", "stockbotpass123", "password used by stock-bot to auth in chatroom")

	flag.Parse()
}

func init() {
	setupFlags()
	tokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log = logger.Sugar()

	chatroom := entity.NewChatroom()
	stockBotConsumer := configStockBotConsumer(chatroom.GetBroadcastChan())

	userInMemory := repository.NewUserInMemory()
	createStockBotUser(userInMemory)

	chatroomService := chatroomsvc.NewService(chatroom, log)
	userService := user.NewService(userInMemory, log)

	go chatroomService.Start()

	r := chi.NewRouter()
	handler.MakePrivateHandlers(r, chatroom, userService, tokenAuth, log)
	handler.MakePublicHandlers(r, userService, tokenAuth)

	log.Info("starting HTTP server on " + addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	if err := stockBotConsumer.Close(); err != nil {
		log.Fatal(err)
	}
}

func configStockBotConsumer(dispatcher chan *entity.Message) rabbitmq.Consumer {
	amqpConn, err := rabbitmq.NewAmqpConnManager(rabbitmq.AmqpConfig{
		User:     rabbitUser,
		Password: rabbitPass,
		Host:     rabbitAddr,
	})
	if err != nil {
		log.Fatal(err)
	}

	consumer := amqpConn.CreateConsumer("stockbot1", TopicChatOutbound)
	if err := consumer.Start(dispatcher); err != nil {
		log.Fatal(err)
	}

	return consumer
}

func createStockBotUser(userRepo user.Repository) {
	if err := userRepo.Create(entity.NewUser("stock-bot", stockBotPassword)); err != nil {
		log.Error(err)
	}
}
