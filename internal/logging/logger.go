package logging

import (
	"log/slog"
	"os"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/config"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func New(cfg *config.Config) *slog.Logger {
	var base slog.Handler

	switch cfg.AppEnv {
	case "production":
		base = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: redactAttr,
		})
	default:
		isTerm := isatty.IsTerminal(os.Stderr.Fd()) ||
			isatty.IsCygwinTerminal(os.Stderr.Fd())

		base = tint.NewHandler(os.Stderr, &tint.Options{
			Level:       slog.LevelDebug,
			ReplaceAttr: redactAttr,
			NoColor:     !isTerm,
		})
	}

	hostname, _ := os.Hostname()
	logger := slog.New(base).With(slog.String("env", cfg.AppEnv), slog.String("hostname", hostname))
	slog.SetDefault(logger)
	return logger
}
