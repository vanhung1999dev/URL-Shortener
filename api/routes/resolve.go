package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/vanhung1999dev/url-shortener/database"
)

func ResolveURL(ctx *fiber.Ctx) error {
	url := ctx.Params("longUrl")

	redisServer0 := database.CreateClient(0)
	defer redisServer0.Close()

	value, err := redisServer0.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		// Key not found in Redis
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found in the database",
		})
	} else if err != nil {
		// Any other actual Redis error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Redis error: " + err.Error(),
		})
	}

	redisServer1 := database.CreateClient(1)
	defer redisServer1.Close()

	redisServer1.Incr(database.Ctx, "counter")

	return ctx.Redirect(value, 301)
}
