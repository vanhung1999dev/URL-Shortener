package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/vanhung1999dev/url-shortener/routes"

	swagger "github.com/gofiber/swagger"
)

// @title       URL Shortener API
// @version     1.0
// @description This is a simple URL shortener written in Go using Fiber
// @host        localhost:3000
// @BasePath    /

func setupRoutes(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/v1/url/:shortID", routes.ResolveURL)
	app.Post("/v1/url", routes.ShortenURL)
}

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println("error", err)
	}

	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen((os.Getenv("APP_PORT"))))
}
