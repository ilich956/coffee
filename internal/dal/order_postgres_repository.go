package dal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"hot-coffee/models"
	myerrors "hot-coffee/internal/myErrors"
)

type postgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{
		db: db,
	}
}

func (r *postgresOrderRepository) GetOrder() ([]models.Order, error) {
	rows, err := r.db.Query("SELECT id, items, total_price, status, created_at, updated_at FROM orders")
	if err != nil {
		slog.Error("failed to get orders", "error", err)
		return nil, myerrors.ErrInternal
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var itemsJSON []byte
		if err := rows.Scan(&order.ID, &itemsJSON, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			slog.Error("failed to scan order", "error", err)
			return nil, myerrors.ErrInternal
		}
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			slog.Error("failed to unmarshal order items", "error", err)
			return nil, myerrors.ErrInternal
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *postgresOrderRepository) GetOrderID(id string) (models.Order, error) {
	var order models.Order
	var itemsJSON []byte
	err := r.db.QueryRow("SELECT id, items, total_price, status, created_at, updated_at FROM orders WHERE id = $1", id).Scan(&order.ID, &itemsJSON, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("order not found", "id", id, "error", err)
			return models.Order{}, myerrors.ErrNotFound
		}
		slog.Error("failed to get order", "id", id, "error", err)
		return models.Order{}, myerrors.ErrInternal
	}
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		slog.Error("failed to unmarshal order items", "error", err)
		return models.Order{}, myerrors.ErrInternal
	}

	return order, nil
}

func (r *postgresOrderRepository) CreateOrder(order models.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		slog.Error("failed to marshal order items", "error", err)
		return myerrors.ErrInternal
	}

	_, err = r.db.Exec("INSERT INTO orders (id, items, total_price, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", order.ID, itemsJSON, order.TotalPrice, order.Status, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		slog.Error("failed to create order", "error", err)
		return myerrors.ErrInternal
	}

	return nil
}

func (r *postgresOrderRepository) UpdateOrder(id string, order models.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		slog.Error("failed to marshal order items", "error", err)
		return myerrors.ErrInternal
	}

	res, err := r.db.Exec("UPDATE orders SET items = $1, total_price = $2, status = $3, updated_at = $4 WHERE id = $5", itemsJSON, order.TotalPrice, order.Status, order.UpdatedAt, id)
	if err != nil {
		slog.Error("failed to update order", "id", id, "error", err)
		return myerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "error", err)
		return myerrors.ErrInternal
	}

	if rowsAffected == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}

func (r *postgresOrderRepository) DeleteOrder(id string) error {
	res, err := r.db.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		slog.Error("failed to delete order", "id", id, "error", err)
		return myerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "error", err)
		return myerrors.ErrInternal
	}

	if rowsAffected == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}

func (r *postgresOrderRepository) CloseOrder(id string) error {
	res, err := r.db.Exec("UPDATE orders SET status = $1 WHERE id = $2", "closed", id)
	if err != nil {
		slog.Error("failed to close order", "id", id, "error", err)
		return myerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "error", err)
		return myerrors.ErrInternal
	}

	if rowsAffected == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
