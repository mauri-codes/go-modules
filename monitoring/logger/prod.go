package logger

import (
	"log/slog"
	"os"
)

func NewProdLogger() *slog.Logger {

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
