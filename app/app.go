package app

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/template/pug"
	"log"
	"os"
	"strconv"
)

var ctx = context.Background()

func Run() {
	app, rdb := boot()

	MountRoot(app, rdb)
	MountAdmin(app, rdb)

	app.Use(func(c *fiber.Ctx) error {
		return c.
			Status(fiber.StatusNotFound).
			SendString("Sorry can't find that!")
	})

	log.Fatal(app.Listen(":3000"))
}

func boot() (*fiber.App, *redis.Client) {
	engine := pug.New("./views", ".pug")
	app := fiber.New(fiber.Config{
		Views:                   engine,
		AppName:                 "luk4s.dev",
		ViewsLayout:             "layouts/application",
		ServerHeader:            "wordprezz",
		EnableTrustedProxyCheck: true,
		Prefork:                 getEnv("PREFORK", "0") == "1",
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(etag.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Accepts("html", "text/html")
		return c.Next()
	})
	app.Static("/public", "./public", fiber.Static{
		Compress: true,
		Browse:   os.Getenv("APP_ENV") == "development",
	})

	redisDb, err := strconv.ParseInt(getEnv("REDIS_DB", "0"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       int(redisDb),
	})

	return app, rdb
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
