package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/vanhung1999dev/url-shortener/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/v1/url/:longUrl", routes.ResolveURL)
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
