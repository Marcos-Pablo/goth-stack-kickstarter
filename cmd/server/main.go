package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/auth"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/home"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMid "github.com/go-chi/chi/v5/middleware"
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

	homeH := home.New(a)
	authH := auth.New(a)

	r := chi.NewRouter()
	r.Use(chiMid.RequestID)
	r.Use(chiMid.RealIP)
	r.Use(chiMid.Recoverer)
	r.Use(slogchi.New(a.Logger))
	r.Use(a.Sessions.LoadAndSave)

	fs := http.FileServer(http.Dir("./assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	r.Route("/auth", func(r chi.Router) {
		r.Get("/sign-in", authH.SignInPage)
		r.Post("/sign-in", authH.SignIn)

		r.Get("/sign-up", authH.SignUpPage)
		r.Post("/sign-up", authH.SignUp)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(a.Sessions))
		r.Get("/", homeH.Index)
		r.Post("/auth/sign-out", authH.SignOut)
	})

	if err = http.ListenAndServe(":8080", r); err != nil {
		a.Logger.Error("failed to initialize server", slog.Any("error", err))
	}
}
