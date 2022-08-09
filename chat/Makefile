help: # Print help on Makefile
	@grep '^[^.#]\+:\s\+.*#' Makefile | \
	sed "s/\(.\+\):\s*\(.*\) #\s*\(.*\)/`printf "\033[93m"`\1`printf "\033[0m"`	\3 [\2]/" | \
	expand -t20

build:
	go build -o ./bin/chatapp ./cmd

start-docker-services:
	@docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

stop-docker-services:
	docker stop rabbitmq || true

run-local:
	@echo "Starting chat app..."
	./bin/chatapp > ../logs/chatapp.log 2>&1 &

stop-local:
	killall chatapp || true