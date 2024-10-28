package router_shortener

import (
	"github.com/gofiber/fiber/v2"
)

func (r ShortenerRoute) ShortenerGetHandler(c *fiber.Ctx) (err error) {
	id := c.Params("id")

	if shortener := r.DB.GetShortener(id); shortener != nil {
		return c.Redirect(shortener.Link)
	}

	return c.Next()
}
