package app

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/config"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/db"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/logging"
)

type App struct {
	Cfg     *config.Config
	DB      *sql.DB
	Queries *db.Queries
	Logger  *slog.Logger
}

func New() (*App, error) {
	cfg, err := config.Load()

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger := logging.New(cfg)

	sqlDB, err := sql.Open("sqlite", buildDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &App{
		Cfg:     cfg,
		DB:      sqlDB,
		Queries: db.New(sqlDB),
		Logger:  logger,
	}, nil
}

func (a *App) Close() error { return a.DB.Close() }

func buildDSN(cfg *config.Config) string {
	return cfg.DbPath +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=foreign_keys(on)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=synchronous(NORMAL)"
}
