# Stock Bot

Authenticate and connect to a chatroom by websocket and listen the commando `/stock=symbol` on chat and responds with the current quote.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Infrastructure](#infrastructure)
    - [RabbitMQ](#rabbitmq)
- [Running Locally](#running-localy)


## Infrastructure
The application needs the following services to run:
- [RabbitMQ](https://www.rabbitmq.com/)

With [Docker](https://www.docker.com/products/docker-desktop/) already installed locally, you can run the following Make command to run all the necessary infraestructure:

```bash
make start-docker-services
```

### RabbitMQ
If you don't have RabbitMQ locally installed, to run a RabbitMQ ina a container run:
```bash
	docker run -d --rm \ 
    --name rabbitmq -p 15672:15672 \
    -p 5672:5672 rabbitmq:3-management
```

## Running Locally

Requirements
- [Golang: 1.18](https://go.dev/dl/)
- [Docker (only to run rabbitmq container)](https://www.docker.com/products/docker-desktop/) 
- [Infrastructure](#infrastructure)

### Using Makefile (recommended)

To run chatapp you can run the following Make command:
```bash
make run-local
```
> This command will run chatapp in background and generate their logs on `../logs` folder

If you want to stop the execution, you can run:
```bash
make stop-local
```

### Using Bash Directly
If you are not a fan of Makefile or can't use it, run the application directly by `go run` in root folder:
```bash
go run ./cmd/main.go
```