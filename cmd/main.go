package main

import (
	"fmt"
	"log"
	"os"

	"hot-coffee/internal/config"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/db"
	"hot-coffee/internal/server"
	"hot-coffee/internal/utils/dir"
	_ "hot-coffee/internal/utils/logger"
)

func main() {
	config.ParseFlags()

	var (
		menuRepo      dal.MenuRepository
		orderRepo     dal.OrderRepository
		inventoryRepo dal.InventoryRepository
	)

	switch *config.Storage {
	case "json":
		dir.CreateDir()
		menuRepo = dal.NewJsonMenuRepository()
		orderRepo = dal.NewJsonOrderRepository()
		inventoryRepo = dal.NewJsonInventoryRepository()
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
		db, err := db.NewDB(dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		menuRepo = dal.NewPostgresMenuRepository(db)
		orderRepo = dal.NewPostgresOrderRepository(db)
		inventoryRepo = dal.NewPostgresInventoryRepository(db)
	}

	server.StartServer(menuRepo, orderRepo, inventoryRepo)
}
