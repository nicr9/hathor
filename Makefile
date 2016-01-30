COMPOSE=docker-compose --project-name hathor

all: build-all hathor-up

debug: build-all
	$(COMPOSE) run backend /bin/bash

build-all:
	$(COMPOSE) build

hathor-up:
	$(COMPOSE) up backend

logs:
	$(COMPOSE) logs
