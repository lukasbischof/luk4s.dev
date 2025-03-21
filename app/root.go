package app

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lukasbischof/luk4s.dev/app/forum"
	"github.com/lukasbischof/luk4s.dev/app/hcaptcha"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
)

func MountRoot(app *fiber.App, db *sql.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return renderIndex(c, db)
	})

	app.Post("/forum", func(c *fiber.Ctx) error {
		return handleForumEntryCreateRequest(c, db)
	})
}

func renderIndex(c *fiber.Ctx, db *sql.DB) error {
	visitorCount, err := IncreaseVisitorCount(db)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	forumEntries, err := GetForumEntries(db)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	return c.Render("index", fiber.Map{
		"forumEntries": forumEntries,
		"visitorCount": message.NewPrinter(language.English).Sprintf("%d", visitorCount),
		"siteKey":      os.Getenv("HCAPTCHA_SITE_KEY"),
	})
}

func handleForumEntryCreateRequest(c *fiber.Ctx, db *sql.DB) error {
	entry := new(forum.Entry)
	if err := c.BodyParser(entry); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	entry.Process()
	if err := entry.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	verify, err := hcaptcha.Verify(entry.CaptchaResponse)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot verify captcha")
	} else if !verify {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid captcha")
	}

	err = SaveForumEntry(db, entry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/")
}
