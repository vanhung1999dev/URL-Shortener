package routes

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/vanhung1999dev/url-shortener/database"
)

func ResolveURL(ctx *fiber.Ctx) error {
	shortID := ctx.Params("shortID")

	redisMain := database.CreateClient(0)
	defer redisMain.Close()

	originalURL, err := redisMain.Get(database.Ctx, shortID).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found in the database",
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Redis error: " + err.Error(),
		})
	}

	redisAnalytics := database.CreateClient(1)
	defer redisAnalytics.Close()

	// Global counter
	redisAnalytics.Incr(database.Ctx, "counter")

	// Optional: per-shortID analytics
	redisAnalytics.Incr(database.Ctx, fmt.Sprintf("clicks:%s", shortID))

	return ctx.Redirect(originalURL, fiber.StatusMovedPermanently)
}
