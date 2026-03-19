package dal

import (
	"database/sql"
	"errors"
	"log/slog"

	"hot-coffee/models"
	myerrors "hot-coffee/internal/myErrors"
)

type postgresInventoryRepository struct {
	db *sql.DB
}

func NewPostgresInventoryRepository(db *sql.DB) InventoryRepository {
	return &postgresInventoryRepository{
		db: db,
	}
}

func (r *postgresInventoryRepository) GetInventory() ([]models.InventoryItem, error) {
	rows, err := r.db.Query("SELECT ingredient_id, name, quantity, unit, threshold FROM inventory")
	if err != nil {
		slog.Error("failed to get inventory", "error", err)
		return nil, myerrors.ErrInternal
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit, &item.Threshold); err != nil {
			slog.Error("failed to scan inventory item", "error", err)
			return nil, myerrors.ErrInternal
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

func (r *postgresInventoryRepository) GetInventoryID(id string) (models.InventoryItem, error) {
	var item models.InventoryItem
	err := r.db.QueryRow("SELECT ingredient_id, name, quantity, unit, threshold FROM inventory WHERE ingredient_id = $1", id).Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit, &item.Threshold)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("inventory item not found", "id", id, "error", err)
			return models.InventoryItem{}, myerrors.ErrNotFound
		}
		slog.Error("failed to get inventory item", "id", id, "error", err)
		return models.InventoryItem{}, myerrors.ErrInternal
	}

	return item, nil
}

func (r *postgresInventoryRepository) CreateInventory(item models.InventoryItem) error {
	_, err := r.db.Exec("INSERT INTO inventory (ingredient_id, name, quantity, unit, threshold) VALUES ($1, $2, $3, $4, $5)", item.IngredientID, item.Name, item.Quantity, item.Unit, item.Threshold)
	if err != nil {
		slog.Error("failed to create inventory item", "error", err)
		return myerrors.ErrInternal
	}

	return nil
}

func (r *postgresInventoryRepository) UpdateInventory(id string, item models.InventoryItem) error {
	res, err := r.db.Exec("UPDATE inventory SET name = $1, quantity = $2, unit = $3, threshold = $4 WHERE ingredient_id = $5", item.Name, item.Quantity, item.Unit, item.Threshold, id)
	if err != nil {
		slog.Error("failed to update inventory item", "id", id, "error", err)
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

func (r *postgresInventoryRepository) DeleteInventory(id string) error {
	res, err := r.db.Exec("DELETE FROM inventory WHERE ingredient_id = $1", id)
	if err != nil {
		slog.Error("failed to delete inventory item", "id", id, "error", err)
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
