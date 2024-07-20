package rest

import (
	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/pkg/internal/database"
)

func initShortener(app *fiber.App, db *database.Database) {
	app.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		if shortener := db.GetShortener(id); shortener != nil {
			return c.Redirect(shortener.Link)
		}
		return c.Next()
	})
}