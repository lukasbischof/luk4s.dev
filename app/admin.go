package app

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"os"
)

func MountAdmin(app *fiber.App, rdb *redis.Client) {
	admin := app.Group("/admin", func(c *fiber.Ctx) error {
		c.Accepts("html", "text/html")
		c.Set("X-Content-Type-Options", "nosniff")
		return c.Next()
	})

	admin.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": os.Getenv("ADMIN_PASSWORD"),
		},
	}))

	admin.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/forum")
	})

	admin.Get("/forum", func(c *fiber.Ctx) error {
		forumEntries, err := GetForumEntries(rdb, ctx)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		return c.Render("admin/forum/index", fiber.Map{
			"forumEntries": forumEntries,
		})
	})

	admin.Post("/forum/:id/delete", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := DeleteForumEntry(rdb, ctx, id); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		return c.Redirect("/admin/forum")
	})
}
