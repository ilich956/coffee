package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"hot-coffee/internal/config"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	_ "hot-coffee/internal/utils/logger"
)

// REPLACE LOGGER

func main() {
	config.ParseFlags()

	// IF DIR EXISTS DONT CREATE JUST PROCEED
	if err := os.Mkdir(*config.Dir, 0o755); err != nil {
		if !os.IsExist(err) {
			slog.Error("Failed to create folder: ", "error", err)
		} else {
			slog.Info("Directory already exists")
		}
	}

	mux := http.NewServeMux()

	orderRepo := dal.NewOrderRepository("orders.json")
	menuRepo := dal.NewMenuRepository("menu_items.json")
	inventoryRepo := dal.NewInventoryRepository("inventory_item.json")

	orderService := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)

	orderHandler := handler.NewOrderHandler(orderService)
	menuHandler := handler.NewMenuHandler(menuService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	// ORDERS
	mux.HandleFunc("GET /orders", orderHandler.HandleGetOrder)
	mux.HandleFunc("GET /orders/{id}", orderHandler.HandleGetOrderID)
	mux.HandleFunc("POST /orders", orderHandler.HandlePostOrder)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.HandlePostOrderClose)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.HandlePutOrderID)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.HandleDeleteOrder)

	// //MENU
	mux.HandleFunc("GET /menu", menuHandler.HandleGetMenu)
	mux.HandleFunc("GET /menu/{id}", menuHandler.HandleGetMenuID)
	mux.HandleFunc("POST /menu", menuHandler.HandlePostMenu)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.HandlePutMenuID)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.HandleDeleteMenuID)

	// //INVENTORY
	mux.HandleFunc("GET /inventory", inventoryHandler.HandleGetInventory)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.HandleGetInventoryID)
	mux.HandleFunc("POST /inventory", inventoryHandler.HandlePostInventory)
	mux.HandleFunc("PUT inventory/{id}", inventoryHandler.HandlePutInventoryID)
	mux.HandleFunc("DELETE inventory/{id}", inventoryHandler.HandleDeleteInventory)

	// //AGREGATIONS
	// mux.HandleFunc("GET /reports/total-sales", handler.)
	// mux.HandleFunc("GET GET /reports/popular-items", handler.)

	if err := http.ListenAndServe(":"+*config.Port, mux); err != nil {
		log.Fatal("Failed to launch server ", err)
	}
}
