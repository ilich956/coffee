package dir

import (
	"hot-coffee/internal/config"
	"log/slog"
	"os"

	myerrors "hot-coffee/internal/myErrors"
)

func CreateDir() {
	if err := os.MkdirAll(*config.Dir, 0o755); err != nil {
		if !os.IsExist(err) {
			slog.Error("Failed to create folder: ", "error", err)
			return
		} else {
			slog.Warn("Directory already exists")
			return
		}
	}
	slog.Info("Directory created", "dir", *config.Dir)
	createJSON()
}

func createJSON() error {
	data := []byte("[]")

	fileNames := []string{"orders", "menu_items", "inventory_item"}

	for _, fileName := range fileNames {
		err := os.WriteFile(*config.Dir+"/"+fileName+".json", data, os.ModePerm)
		if err != nil {
			slog.Error("Failed to write file", "error", err, "file name", fileName)
			return myerrors.ErrFailWrite
		}

		slog.Info("File created", "file", fileName)
	}

	return nil
}
