package app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/template/pug"
	"github.com/lukasbischof/luk4s.dev/app/forum"
	"github.com/qinains/fastergoding"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"strconv"
)

var ctx = context.Background()

func Run() {
	if os.Getenv("APP_ENV") == "development" {
		fastergoding.Run()
		fmt.Println("Starting with Fastergoding")
	}

	app, rdb := boot()

	app.Get("/", func(c *fiber.Ctx) error {
		return renderIndex(c, rdb)
	})

	app.Post("/forum", func(c *fiber.Ctx) error {
		return renderForum(c, rdb)
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

	visitorCount, err := IncreaseVisitorCount(rdb, ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	forumEntries, err := GetForumEntries(rdb, ctx)

	return c.Render("index", fiber.Map{
		"forumEntries": forumEntries,
		"ip":           c.IP(),
		"visitorCount": message.NewPrinter(language.English).Sprintf("%d", visitorCount),
	})
}

func renderForum(c *fiber.Ctx, rdb *redis.Client) error {
	c.Accepts("html", "text/html")

	entry := new(forum.Entry)
	if err := c.BodyParser(entry); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	entry.Process()
	if err := entry.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := SaveForumEntry(rdb, ctx, entry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/")
}

func boot() (*fiber.App, *redis.Client) {
	engine := pug.New("./views", ".pug")
	app := fiber.New(fiber.Config{
		Views:                   engine,
		AppName:                 "luk4s.dev",
		ViewsLayout:             "layouts/application",
		ServerHeader:            "wordprezz",
		EnableTrustedProxyCheck: true,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(etag.New())

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
