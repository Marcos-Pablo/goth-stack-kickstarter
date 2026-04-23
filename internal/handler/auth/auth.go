package auth

import (
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
)

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{app: app}
}
