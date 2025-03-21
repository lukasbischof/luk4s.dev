package app

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"os"
	"strconv"
)

func MountAdmin(app *fiber.App, db *sql.DB) {
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
		forumEntries, err := GetForumEntries(db)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		return c.Render("admin/forum/index", fiber.Map{
			"forumEntries": forumEntries,
		})
	})

	admin.Post("/forum/:id/delete", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		if err = DeleteForumEntry(db, id); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}

		return c.Redirect("/admin/forum")
	})
}
