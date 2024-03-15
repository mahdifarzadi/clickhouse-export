.PHONY: install
install:
	go install github.com/mahdifarzadi/clickhouse-export@latest

.PHONY: build
build:
	GOFLAGS=-mod=mod go build -o bin/clickhouse-export main.go

.PHONY: run
run: build
	./bin/clickhouse-export

.PHONY: dev
dev: build
	./bin/clickhouse-export -c clickhouse-export.yaml