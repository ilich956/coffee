all: build
deploy: gofumpt build

gofumpt:
	gofumpt -l -w .

build:
	go build -o coffee ./cmd/main.go