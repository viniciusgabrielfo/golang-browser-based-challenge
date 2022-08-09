help: # Print help on Makefile
	@grep '^[^.#]\+:\s\+.*#' Makefile | \
	sed "s/\(.\+\):\s*\(.*\) #\s*\(.*\)/`printf "\033[93m"`\1`printf "\033[0m"`	\3 [\2]/" | \
	expand -t20

build:
	go build -o ./bin/stock-bot ./cmd

start-docker-services:
	@docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

stop-docker-services:
	docker stop rabbitmq || true

run-local:
	@echo "Starting stock-bot app..."
	./bin/stock-bot > ../logs/stock-bot.log 2>&1 &

stop-local:
	killall stock-bot || true