package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/pug"
	"github.com/qinains/fastergoding"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"strconv"
)

var ctx = context.Background()

func main() {
	if os.Getenv("APP_ENV") == "development" {
		fastergoding.Run()
		fmt.Println("Starting with Fastergoding")
	}

	app, rdb := boot()

	app.Get("/", func(c *fiber.Ctx) error {
		return renderIndex(c, rdb)
	})

	app.Static("/public", "./public", fiber.Static{
		Compress: true,
		Browse:   true,
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.
			Status(fiber.StatusNotFound).
			SendString("Sorry can't find that!")
	})

	log.Fatal(app.Listen(":3000"))
}

func renderIndex(c *fiber.Ctx, rdb *redis.Client) error {
	c.Accepts("html", "text/html")

	visitorCount, err := IncreaseVisitorCount(rdb)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	return c.Render("index", fiber.Map{
		"visitorCount": message.NewPrinter(language.English).Sprintf("%d", visitorCount),
	})
}

func boot() (*fiber.App, *redis.Client) {
	engine := pug.New("./views", ".pug")
	app := fiber.New(fiber.Config{
		Views:        engine,
		AppName:      "luk4s.dev",
		ViewsLayout:  "layouts/application",
		ServerHeader: "Luk4serve",
		ETag:         true,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

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
