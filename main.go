package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
	"github.com/qinains/fastergoding"
	"log"
)

func main() {
	fastergoding.Run()
	engine := pug.New("./views", ".pug")
	app := fiber.New(fiber.Config{
		Views:       engine,
		AppName:     "luk4s.dev",
		ViewsLayout: "layouts/application",
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"test": "Testin",
		})
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
