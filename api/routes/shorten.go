package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/sony/sonyflake"
	"github.com/vanhung1999dev/url-shortener/database"
	"github.com/vanhung1999dev/url-shortener/helpers"
)

type request struct {
	LongURL string `json:"long_url"`
}

type response struct {
	ShortURL string `json:"short_url"`
}

func ShortenURL(ctx *fiber.Ctx) error {
	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if !govalidator.IsURL(body.LongURL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	body.LongURL = helpers.EnforceHTTP(body.LongURL)

	// Rate limit per second
	redisRate := database.CreateClient(1)
	defer redisRate.Close()

	rateKey := "rate:" + ctx.IP()
	exists, _ := redisRate.Exists(database.Ctx, rateKey).Result()
	if exists > 0 {
		return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Rate limit exceeded, try again in 1 second"})
	}
	redisRate.Set(database.Ctx, rateKey, "1", 1*time.Second)

	// Check if long URL already exists
	redisMain := database.CreateClient(0)
	defer redisMain.Close()

	existingShortKey := "url_map:" + body.LongURL
	existingID, _ := redisMain.Get(database.Ctx, existingShortKey).Result()
	if existingID != "" {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "URL already shortened"})
	}

	// Generate Snowflake ID
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return 1, nil
		},
	})

	if sf == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to initialize ID generator"})
	}

	id, err := sf.NextID()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ID generation error"})
	}

	shortID := strconv.FormatUint(id, 36) // base36 encoding

	// Save shortID → longURL and longURL → shortID
	if err := redisMain.Set(database.Ctx, shortID, body.LongURL, 24*time.Hour).Err(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store URL"})
	}

	// Map longURL → shortID to detect duplicates
	_ = redisMain.Set(database.Ctx, existingShortKey, shortID, 24*time.Hour).Err()

	// Return response
	return ctx.Status(fiber.StatusOK).JSON(response{
		ShortURL: os.Getenv("DOMAIN") + "/" + shortID,
	})
}
