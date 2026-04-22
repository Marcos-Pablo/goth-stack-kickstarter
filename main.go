package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/views"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Static("/assets", "./assets")

	router.GET("/", func(ctx *gin.Context) {
		render(ctx, http.StatusOK, views.Home())
	})

	log.Fatal(router.Run(":8080"))
}

func render(ctx *gin.Context, status int, cmp templ.Component) {
	ctx.Status(status)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := cmp.Render(context.Background(), ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to render template",
		})
	}
}
