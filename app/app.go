package app

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/template/pug"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var ctx = context.Background()

func Run() {
	app, db := boot()

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	MountRoot(app, db)
	MountAdmin(app, db)

	app.Use(func(c *fiber.Ctx) error {
		return c.
			Status(fiber.StatusNotFound).
			SendString("Sorry can't find that!")
	})

	log.Fatal(app.Listen(":3000"))
}

func boot() (*fiber.App, *sql.DB) {
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

	db, err := sql.Open("sqlite3", getEnv("APP_DB", "./luk4s.db"))
	if err != nil {
		log.Fatal(err)
	}

	return app, db
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
