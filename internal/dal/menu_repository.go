package dal

import (
	"encoding/json"
	"hot-coffee/internal/config"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"log/slog"
	"os"

	myerrors "hot-coffee/internal/myErrors"
)

type MenuRepository interface {
	GetMenu() ([]models.MenuItem, error)
	GetMenuID(id string) (models.MenuItem, error)
	CreateMenu(newMenuItem models.MenuItem) error
	UpdateMenu(id string, newMenu models.MenuItem) error
	DeleteMenu(id string) error
}

type jsonMenuRepository struct {
	filepath string
}

func NewMenuRepository(filepath string) MenuRepository {
	return &jsonMenuRepository{filepath: filepath}
}

func (o *jsonMenuRepository) GetMenu() ([]models.MenuItem, error) {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return []models.MenuItem{}, myerrors.ErrFailOpenJson
	}

	var menuItems []models.MenuItem

	if err := json.Unmarshal(byteValue, &menuItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return []models.MenuItem{}, myerrors.ErrFailUnmarshal
	}

	return menuItems, nil
}

func (m *jsonMenuRepository) GetMenuID(id string) (models.MenuItem, error) {
	byteValue, err := utils.ReadFile(m.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", m.filepath)
		return models.MenuItem{}, myerrors.ErrFailOpenJson
	}

	var menus []models.MenuItem
	if err := json.Unmarshal(byteValue, &menus); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return models.MenuItem{}, myerrors.ErrFailUnmarshal
	}

	for _, menu := range menus {
		if menu.ID == id {
			return menu, nil
		}
	}

	return models.MenuItem{}, myerrors.ErrNotFound
}

func (o *jsonMenuRepository) CreateMenu(newMenuItem models.MenuItem) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)
		return myerrors.ErrFailOpenJson
	}

	var menuItems []models.MenuItem
	if err := json.Unmarshal(byteValue, &menuItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	menuItems = append(menuItems, newMenuItem)

	filestring, _ := json.MarshalIndent(menuItems, "", "  ")
	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}

func (m *jsonMenuRepository) UpdateMenu(id string, newMenu models.MenuItem) error {
	byteValue, err := utils.ReadFile(m.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", m.filepath)
		return myerrors.ErrFailUnmarshal
	}

	var menus []models.MenuItem
	if err := json.Unmarshal(byteValue, &menus); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	var isFound bool

	for i := range menus {
		if menus[i].ID == id {
			menus[i] = newMenu
			isFound = true
		}
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}

	filestring, _ := json.MarshalIndent(menus, "", "  ")
	os.WriteFile(*config.Dir+"/"+m.filepath, filestring, os.ModePerm)

	return nil
}

func (o *jsonMenuRepository) DeleteMenu(id string) error {
	byteValue, err := utils.ReadFile(o.filepath)
	if err != nil {
		slog.Error("Failed to open", "error", err, "file path", o.filepath)

		return myerrors.ErrFailOpenJson
	}

	var menuItems []models.MenuItem

	if err := json.Unmarshal(byteValue, &menuItems); err != nil {
		slog.Error("Failed to unmarshal", "error", err)

		return myerrors.ErrFailUnmarshal
	}

	var newMenuItems []models.MenuItem
	var isFound bool

	for i := range menuItems {
		if menuItems[i].ID == id {
			isFound = true
			continue
		}
		newMenuItems = append(newMenuItems, menuItems[i])
	}

	if !isFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
		return myerrors.ErrNotFound
	}

	filestring, err := json.MarshalIndent(newMenuItems, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal", "error", err)

		return myerrors.ErrFailMarshal
	}

	os.WriteFile(*config.Dir+"/"+o.filepath, filestring, os.ModePerm)

	return nil
}
