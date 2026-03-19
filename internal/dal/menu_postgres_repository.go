package dal

import (
	"database/sql"
	"errors"
	"log/slog"

	myerrors "hot-coffee/internal/myErrors"
	"hot-coffee/models"
)

type postgresMenuRepository struct {
	db *sql.DB
}

func NewPostgresMenuRepository(db *sql.DB) MenuRepository {
	return &postgresMenuRepository{
		db: db,
	}
}

func (r *postgresMenuRepository) GetMenu() ([]models.MenuItem, error) {
	rows, err := r.db.Query("SELECT id, name, price, description FROM menu")
	if err != nil {
		slog.Error("failed to get menu", "error", err)
		return nil, myerrors.ErrInternal
	}
	defer rows.Close()

	var menu []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description); err != nil {
			slog.Error("failed to scan menu item", "error", err)
			return nil, myerrors.ErrInternal
		}
		menu = append(menu, item)
	}

	return menu, nil
}

func (r *postgresMenuRepository) GetMenuID(id string) (models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.QueryRow("SELECT id, name, price, description FROM menu WHERE id = $1", id).Scan(&item.ID, &item.Name, &item.Price, &item.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("menu item not found", "id", id, "error", err)
			return models.MenuItem{}, myerrors.ErrNotFound
		}
		slog.Error("failed to get menu item", "id", id, "error", err)
		return models.MenuItem{}, myerrors.ErrInternal
	}

	return item, nil
}

func (r *postgresMenuRepository) CreateMenu(item models.MenuItem) error {
	_, err := r.db.Exec("INSERT INTO menu (id, name, price, description) VALUES ($1, $2, $3, $4)", item.ID, item.Name, item.Price, item.Description)
	if err != nil {
		slog.Error("failed to create menu item", "error", err)
		return myerrors.ErrInternal
	}

	return nil
}

func (r *postgresMenuRepository) UpdateMenu(id string, item models.MenuItem) error {
	res, err := r.db.Exec("UPDATE menu SET name = $1, price = $2, description = $3 WHERE id = $4", item.Name, item.Price, item.Description, id)
	if err != nil {
		slog.Error("failed to update menu item", "id", id, "error", err)
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

func (r *postgresMenuRepository) DeleteMenu(id string) error {
	res, err := r.db.Exec("DELETE FROM menu WHERE id = $1", id)
	if err != nil {
		slog.Error("failed to delete menu item", "id", id, "error", err)
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
