package main

import (
	"log"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/auth"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/home"
	"github.com/gin-gonic/gin"

	_ "modernc.org/sqlite"
)

func main() {
	app, err := app.New()

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	_ = auth.New(app)
	homeH := home.New(app)

	engine := gin.Default()
	engine.Static("/assets", "./assets")

	engine.GET("/", homeH.Index)

	log.Fatal(engine.Run(":8080"))
}
