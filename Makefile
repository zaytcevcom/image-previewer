BIN := "./bin/previewer"
DOCKER_IMG="previewer:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./api

init: docker-down-clear \
	clear \
	docker-network \
	docker-pull docker-build docker-up

up: docker-network docker-up
down: docker-down
restart: down up

docker-pull:
	docker compose -f ./deployments/development/docker-compose.yml pull

docker-build:
	docker compose -f ./deployments/development/docker-compose.yml build --pull

docker-up:
	docker compose -f ./deployments/development/docker-compose.yml up -d

docker-down:
	docker compose -f ./deployments/development/docker-compose.yml down --remove-orphans

docker-down-clear:
	docker compose -f ./deployments/development/docker-compose.yml down -v --remove-orphans

docker-network:
	docker network create previewer_network || true

clear:
	rm -rf var/postgres/* var/rabbitmq/*

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/previewer

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version


test:
	go test -race -count 100 -timeout 60m ./internal/...


remove-lint-deps:
	rm $(which golangci-lint)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

integration-tests:
	set -e ;\
	docker compose -f ./deployments/test/docker-compose.yml up --build -d ;\
	test_status_code=0 ;\
	docker compose -f ./deployments/test/docker-compose.yml run integration_tests go test || test_status_code=$$? ;\
	docker compose -f ./deployments/test/docker-compose.yml down \
                --rmi local \
        		--volumes \
        		--remove-orphans \
        		--timeout 60; \
	echo $$test_status_code ;\
	exit $$test_status_code ;