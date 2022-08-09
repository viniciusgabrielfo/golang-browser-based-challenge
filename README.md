# Golang Browser-Based Chat Application

This is a browser-based chat application written in Golang. In general, you can talk in a chatroom with others authenticated users on `/room`.

The "database" users is managed in memory and you can register new users in `/signup` page.

This is a repository with 2 modules/services:

## [ChatApp](./chat/)
Responsible to manage users and the chatroom.

## [Stock-Bot](./stock-bot/)
Responsible to handle the command `/stock=symbol` on the chatroom. Works as if you were a user within the chat.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [How to Use](#how-to-use)
- [Infrastructure](#infrastructure)
    - [RabbitMQ](#rabbitmq)
- [Running Locally](#running-localy)
- [ChatApp Endpoints](#Endpoints)

## How to Use
**It will only work if all services are running, If not please visit [Infrastructure](#infrastructure) and [Running Locally](#running-localy)**

1. Access `localhost:8000/login`
2. If you don't have a user, click to Signup or access `localhost:8000/signup`
3. Login with your user and password (you will be redirected to chatroom, if not access `localhost:8000/chatroom`) 

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
To run any Make or Go command at the project root **is necessary a go.work file**. To generate run:
```bash
make generate-local-workspace
```

To run chatapp and stock-bot together you can run the following Make command:
```bash
make run-services-local
```
> This command will run services in background and generate their logs on `/logs` folder

If you want to stop the execution, you can run:
```bash
make stop-services-local
```

### Using Bash Directly
If you are not a fan of Makefile or can't use it, run the application directly by `go run` in root folder:
```bash
go work init
go work use -r . 
go run chat/cmd/main.go
go run stock-bot/cmd/main.go
```
> It's important to generate a go.work to get a better experience with multi-modules

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