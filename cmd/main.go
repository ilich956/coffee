package main

import (
	"hot-coffee/internal/config"
	"hot-coffee/internal/server"
	"hot-coffee/internal/utils/dir"
	_ "hot-coffee/internal/utils/logger"
)

func main() {
	config.ParseFlags()

	dir.CreateDir()

	server.StartServer()
}
