package home

import (
	"net/http"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/middleware"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/views"
)

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.UserFrom(r.Context())
	handler.Render(w, r, http.StatusOK, views.Home(views.User{
		Email: user.Email,
		Name:  user.Name,
	}))
}
