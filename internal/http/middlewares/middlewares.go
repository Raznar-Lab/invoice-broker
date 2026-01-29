package middlewares

import (
	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
)

type Middlewares struct {
	Logger gin.HandlerFunc
}

func New(c *configs.Config) *Middlewares {
	return &Middlewares{
		Logger: LogMiddleware,
	}
}
