package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func Pretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling: %v", err)
	}
	return string(b)
}

func InfoPretty(ctx context.Context, log *slog.Logger, msg string, v any) {
	if !log.Enabled(ctx, slog.LevelInfo) {
		return
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Info("dump failed", "msg", msg, "err", err)
		return
	}
	log.InfoContext(ctx, msg)
	_, _ = os.Stdout.Write(b)
	_, _ = os.Stdout.Write([]byte("\n"))
}

func DebugPretty(ctx context.Context, log *slog.Logger, msg string, v any) {
	if !log.Enabled(ctx, slog.LevelDebug) {
		return
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Info("dump failed", "msg", msg, "err", err)
		return
	}
	log.InfoContext(ctx, msg)
	_, _ = os.Stdout.Write(b)
	_, _ = os.Stdout.Write([]byte("\n"))
}

func ErrorPretty(ctx context.Context, log *slog.Logger, msg string, v any) {
	if !log.Enabled(ctx, slog.LevelError) {
		return
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Info("dump failed", "msg", msg, "err", err)
		return
	}
	log.InfoContext(ctx, msg)
	_, _ = os.Stdout.Write(b)
	_, _ = os.Stdout.Write([]byte("\n"))
}
