package home

import (
	"net/http"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/views"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Index(ctx *gin.Context) {
	handler.Render(ctx, http.StatusOK, views.Home())
}
