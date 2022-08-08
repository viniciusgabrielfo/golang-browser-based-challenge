help: # Print help on Makefile
	@grep '^[^.#]\+:\s\+.*#' Makefile | \
	sed "s/\(.\+\):\s*\(.*\) #\s*\(.*\)/`printf "\033[93m"`\1`printf "\033[0m"`	\3 [\2]/" | \
	expand -t20

build: # builds all services
	@echo "Building all services..."
	@go build -o ./chat/bin/chatapp ./chat/cmd
	@go build -o ./stock-bot/bin/stock-bot ./stock-bot/cmd
	@echo "All services were built. Check ./bin in each module folder"

start-chat-app:
	@go run chat/cmd/main.go

start-stock-bot:
	@go run stock-bot/cmd/main.go

start-docker-services:
	@docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

stop-docker-services:
	docker stop rabbitmq || true

run-services-local: build
	@echo "Starting chat app..."
	./chat/bin/chatapp > logs/chatapp.log 2>&1 &
	@sleep 2
	@echo "Starting stock-bot app..."
	./stock-bot/bin/stock-bot > logs/stock-bot.log 2>&1 &

stop-services-local:
	killall stock-bot || true
	killall chatapp || true

generate-local-workspace: # generates a go workspace file for the project (easily to work with multi-modules locally)
	go work init
	go work use -r . 