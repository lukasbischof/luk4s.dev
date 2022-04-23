package app

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/lukasbischof/luk4s.dev/app/forum"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
)

func MountRoot(app *fiber.App, rdb *redis.Client) {
	app.Get("/", func(c *fiber.Ctx) error {
		return renderIndex(c, rdb)
	})

	app.Post("/forum", func(c *fiber.Ctx) error {
		return renderForum(c, rdb)
	})
}

func renderIndex(c *fiber.Ctx, rdb *redis.Client) error {
	visitorCount, err := IncreaseVisitorCount(rdb, ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	forumEntries, err := GetForumEntries(rdb, ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	return c.Render("index", fiber.Map{
		"forumEntries": forumEntries,
		"visitorCount": message.NewPrinter(language.English).Sprintf("%d", visitorCount),
	})
}

func renderForum(c *fiber.Ctx, rdb *redis.Client) error {
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
