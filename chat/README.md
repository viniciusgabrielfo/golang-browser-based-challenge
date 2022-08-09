# Chat App

Responsible to manage users and the chatroom. Listen stock-bot by RabbitMQ queue `chat.outbound`.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Infrastructure](#infrastructure)
    - [RabbitMQ](#rabbitmq)
- [Running Locally](#running-localy)
- [ChatApp Endpoints](#Endpoints)


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

## Endpoints
| Endpoint | Description | Type | Need JWT |
|---       |---          |---   |---   
| GET /login | Render the login page | Render Page |  |
| GET /signup | Render the signup page | Render Page | |
| GET /room | Render the chatroom page | Render Page | :white_check_mark: 
| POST /auth | Authenticate with the server | API |  |
| POST /register | Register a new user on the server | API | |
| GET /ws | Connect to chatroom websocket | API | :white_check_mark: |

### POST /auth
Payload:
```json
{
    "nick": "nick",
    "password": "pass"
}
```

### POST /register
Payload:
```json
{
    "nick": "nick",
    "password": "pass"
}
```