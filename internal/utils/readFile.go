package utils

import (
	"io"
	"log/slog"
	"os"

	"hot-coffee/internal/config"
)

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(*config.Dir + "/" + filePath)
	if err != nil {
		slog.Error("Failed to open", "error", err)
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read", "error", err)
		return nil, err
	}

	return byteValue, nil
}
