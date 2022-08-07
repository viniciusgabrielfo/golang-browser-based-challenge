package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log = logger.Sugar()
	log.Info("starting chat app...")

	chatroom := entity.NewChatroom()

	amqpConn := createAmqpConnManager()
	stockBotConsumer := amqpConn.CreateConsumer("chatroom", TopicChatOutbound)
	if err := stockBotConsumer.Start(chatroom.GetBroadcastChan()); err != nil {
		log.Fatal(err)
	}

	httpServer := startHTTPServer(chatroom)

	<-stop
	log.Info("starting graceful shutdown...")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}

	if err := stockBotConsumer.Close(); err != nil {
		log.Fatal(err)
	}

	if err := amqpConn.Close(); err != nil {
		log.Fatal(err)
	}

	log.Info("app finished")
}

func createAmqpConnManager() *rabbitmq.AmqpConnManager {
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

	if err := amqpConn.QueueDeclare(TopicChatOutbound); err != nil {
		log.Fatal(err)
	}

	if err := amqpConn.QeueBind(TopicChatOutbound, rabbitRoutingKey, ExchangeChat); err != nil {
		log.Fatal(err)
	}

	return amqpConn
}

func startHTTPServer(chatroom *entity.Chatroom) *http.Server {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	}))

	userInMemory := repository.NewUserInMemory()
	createStockBotUser(userInMemory)

	chatroomService := chatroomsvc.NewService(chatroom, log)
	userService := user.NewService(userInMemory, log)

	go chatroomService.Start()

	handler.MakePrivateHandlers(r, chatroom, userService, tokenAuth, log)
	handler.MakePublicHandlers(r, userService, tokenAuth)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info("starting HTTP server on " + addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	return server
}

func createStockBotUser(userRepo user.Repository) {
	if err := userRepo.Create(entity.NewUser("stock-bot", stockBotPassword)); err != nil {
		log.Error(err)
	}
}
