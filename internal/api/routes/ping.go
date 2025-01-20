package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
)

func PingRoute(customRouter *CustomRouter) {
	customRouter.routerGroup.GET("/", handlers.PingHandler)
}
