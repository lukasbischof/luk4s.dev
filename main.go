package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/pug"
	"github.com/qinains/fastergoding"
	"log"
	"os"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		fastergoding.Run()
		fmt.Println("Starting with Fastergoding")
	}

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

	app.Get("/", func(ctx *fiber.Ctx) error {
		ctx.Accepts("html", "text/html")
		return ctx.Render("index", fiber.Map{})
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
