BIN_APP := "./bin/calendar"
BIN_MIGRATOR := "./bin/migrator"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# Build
build-app: 
	go build -v -o $(BIN_APP) -ldflags "$(LDFLAGS)" ./cmd/calendar
build-migrator:
	go build -v -o ${BIN_MIGRATOR} ./cmd/migrator

# Run
run-app: build
	$(BIN_APP) -config ./configs/config.json
run-migrator: build
	$(BIN_MIGRATOR) -config ./configs/config.json

# Docker
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_APP) version

# Checks
test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

# Code generation
generate:
	rm -rf internal/servers/grpc/generated
	mkdir -p internal/servers/grpc/generated

	protoc \
		--proto_path=api/event_service \
		--go_out=internal/servers/grpc/generated \
		--go-grpc_out=internal/servers/grpc/generated \
		api/event_service/*.proto

.PHONY: build run build-img run-img version test lint
