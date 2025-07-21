package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vanhung1999dev/url-shortener/database"
	"github.com/vanhung1999dev/url-shortener/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"reate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(ctx *fiber.Ctx) error {
	body := new(request)

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// implement rate limiting
	redisServer1 := database.CreateClient(1)
	defer redisServer1.Close()

	val, err := redisServer1.Get(database.Ctx, ctx.IP()).Result()

	if err == redis.Nil {
		_ = redisServer1.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			limit, _ := redisServer1.TTL(database.Ctx, ctx.IP()).Result()

			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceeded", "rate_limit_reset": limit / time.Nanosecond / time.Minute})
		}

	}

	// check if the input if an actual URL
	if !govalidator.IsURL(body.URL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Invald Domain"})
	}

	// enforce HTTPS, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	redisServer0 := database.CreateClient(0)
	defer redisServer0.Close()

	value, _ := redisServer0.Get(database.Ctx, id).Result()

	if value != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Url custom short is already use"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = redisServer0.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != redis.Nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to connect to server"})
	}

	res := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	redisServer1.Decr(database.Ctx, ctx.IP())

	val, _ = redisServer1.Get(database.Ctx, ctx.IP()).Result()
	res.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := redisServer1.TTL(database.Ctx, ctx.IP()).Result()
	res.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	res.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(res)
}
