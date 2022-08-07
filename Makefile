start-docker-services:
	docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

stop-docker-services:
	docker stop rabbitmq || true