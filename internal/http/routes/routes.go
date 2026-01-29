package routes

import (
	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/http/middlewares"
	paypal_route "raznar.id/invoice-broker/internal/http/routes/paypal"
	xendit_route "raznar.id/invoice-broker/internal/http/routes/xendit"
	"raznar.id/invoice-broker/internal/services"
)

func Init(r *gin.Engine, conf *configs.Config, s *services.Services) {
	mw := middlewares.New(conf)
	// Global Middlewares
	r.Use(mw.Logger) // Your zerolog middleware
	r.Use(gin.Recovery())

	// 2. Gateway Routes Group
	api := r.Group("/api")

	// Initialize the Xendit Sub-Router
	xenditRouter := xendit_route.New(conf, s, mw, api)
	paypalRouter := paypal_route.New(conf, s, mw, api)
	xenditRouter.Register()
	paypalRouter.Register()

}
