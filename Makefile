.PHONY: all build clean deploy gofumpt run run-postgres docker-up docker-down docker-build docker-logs

all: build

build:
	go build -o coffee ./cmd/main.go

clean:
	rm -f coffee

deploy: gofumpt build

gofumpt:
	gofumpt -l -w .

run: build
	./coffee

run-postgres: build
	./coffee --storage=postgres

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

docker-logs:
	docker-compose logs -f
