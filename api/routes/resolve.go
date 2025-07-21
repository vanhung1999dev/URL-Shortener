package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/vanhung1999dev/url-shortener/database"
)

func ResolveURL(ctx *fiber.Ctx) error {
	url := ctx.Params("url")

	redisServer0 := database.CreateClient(0)
	defer redisServer0.Close()

	value, err := redisServer0.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found in the database"})
	} else if err != redis.Nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot connect to database"})
	}

	redisServer1 := database.CreateClient(1)
	defer redisServer1.Close()

	redisServer1.Incr(database.Ctx, "counter")

	return ctx.Redirect(value, 301)
}
