.PHONY: help
help:
	@echo 'Management commands:'
	@echo
	@echo 'Usage:'
	@echo '    make help            Print this help message.'
	@echo '    make all             Lint and build.'
	@echo '    make build           Compile project.'
	@echo '    make clean           Clean directory tree.'
	@echo '    make debug           Run with debugger.'
	@echo '    make dev             Run without compiling.'
	@echo '    make fix             Fix small linting problems.'
	@echo '    make lint            Run static analysis on source code.'
	@echo '    make test            Run tests.'
	@echo

.PHONY: all
all: clean lint build

.PHONY: b
b: build

.PHONY: build
build:
	go build cmd/main.go

.PHONY: clean
clean:

.PHONY: debug
debug:
	dlv debug cmd/main.go

.PHONY: dev
dev:
	go run cmd/main.go

.PHONY: fix
fix:
	golangci-lint run --fix

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test
