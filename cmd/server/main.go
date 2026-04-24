package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/auth"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/home"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samber/slog-chi"

	_ "modernc.org/sqlite"
)

func main() {
	a, err := app.New()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	defer a.Close()

	_ = auth.New(a)
	homeH := home.New(a)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(slogchi.NewWithConfig(a.Logger, slogchi.Config{
		WithRequestID: true,
	}))

	fs := http.FileServer(http.Dir("./assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	r.Get("/", homeH.Index)

	if err = http.ListenAndServe(":8080", r); err != nil {
		a.Logger.Error("failed to initialize server", slog.Any("error", err))
	}
}
