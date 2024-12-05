package service

import (
	"encoding/json"
	"log/slog"

	"hot-coffee/internal/dal"
	myerrors "hot-coffee/internal/myErrors"
	"hot-coffee/internal/utils/validation"
	"hot-coffee/models"
)

type MenuService interface {
	ServiceGetMenu() ([]byte, error)
	ServiceGetMenuID(id string) ([]byte, error)
	ServiceCreateMenu(newMenuItem []byte) error
	ServiceUpdateMenu(id string, newMenu []byte) error
	ServiceDeleteMenu(id string) error
}

type menuService struct {
	menuRepo dal.MenuRepository
}

func NewMenuService(repo dal.MenuRepository) MenuService {
	return &menuService{menuRepo: repo}
}

func (m *menuService) ServiceCreateMenu(newMenuItem []byte) error {
	var menu models.MenuItem
	if err := json.Unmarshal(newMenuItem, &menu); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	if err := validation.CheckMenu(menu); err != nil {
		return err
	}

	checkMenuID, _ := m.menuRepo.GetMenuID(menu.ID)
	if checkMenuID.ID == menu.ID {
		slog.Error("Failed to create menu", "error", myerrors.ErrIDExist)
		return myerrors.ErrIDExist
	}

	return m.menuRepo.CreateMenu(menu)
}

func (m *menuService) ServiceDeleteMenu(id string) error {
	return m.menuRepo.DeleteMenu(id)
}

func (m *menuService) ServiceGetMenu() ([]byte, error) {
	menu, err := m.menuRepo.GetMenu()
	if err != nil {
		return nil, err
	}

	jsonFile, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		return nil, myerrors.ErrFailMarshal
	}

	return jsonFile, nil
}

func (m *menuService) ServiceGetMenuID(id string) ([]byte, error) {
	menu, err := m.menuRepo.GetMenuID(id)
	if err == myerrors.ErrNotFound {
		slog.Error("Failed to find", "error", myerrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	jsonFile, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		return nil, myerrors.ErrFailMarshal
	}
	return jsonFile, nil
}

func (m *menuService) ServiceUpdateMenu(id string, newMenu []byte) error {
	var menu models.MenuItem
	if err := json.Unmarshal(newMenu, &menu); err != nil {
		slog.Error("Failed to unmarshal", "error", err)
		return myerrors.ErrFailUnmarshal
	}

	if err := validation.CheckMenu(menu); err != nil {
		return err
	}

	return m.menuRepo.UpdateMenu(id, menu)
}
