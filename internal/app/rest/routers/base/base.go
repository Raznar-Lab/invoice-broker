package base_router

import (
	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/pkg/database"
)

type BaseRouter struct {
	App    *fiber.App
	Config *configs.Config
	DB     *database.Database
}

func (r *BaseRouter) Set(app *fiber.App, config *configs.Config, db *database.Database) {
	r.App = app
	r.Config = config
	r.DB = db
}

type IBaseRouter interface {
	Init(g fiber.Router)
}
