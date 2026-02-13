package logger

import (
	"log/slog"
	"os"
)

func NewDevLogger() *slog.Logger {

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,

		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey, slog.LevelKey:
				return slog.Attr{}
			}
			return a
		},
	})

	return slog.New(handler)
}
