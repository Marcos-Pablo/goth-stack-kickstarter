package handler

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(ctx *gin.Context, status int, cmp templ.Component) {
	ctx.Status(status)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := cmp.Render(context.Background(), ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to render template",
		})
	}
}
