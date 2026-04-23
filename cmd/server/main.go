package main

import (
	"log"
	"net/http"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/auth"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/home"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "modernc.org/sqlite"
)

func main() {
	app, err := app.New()

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	_ = auth.New(app)
	homeH := home.New(app)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("./assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	r.Get("/", homeH.Index)

	log.Fatal(http.ListenAndServe(":8080", r))
}
