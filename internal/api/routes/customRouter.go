package routes

import (
	"bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/internal/config"
	"bizarre-vpn-api/internal/storage/sqlite"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type CustomRouter struct {
	_router       *gin.Engine
	routerGroup   *gin.RouterGroup
	log           *slog.Logger
	cfg           *config.Config
	storage       *sqlite.Storage
	botSharedData *bot.BotSharedData
}

func (cr *CustomRouter) AddGroup(relativePath string, groupHandler func(customRouter *CustomRouter)) {
	if relativePath[0] != '/' {
		panic(fmt.Sprintf(
			"you need to provide correctly paths, basePath: %v, currentPath: %v",
			cr.routerGroup.BasePath(),
			relativePath,
		))
	}

	newGroup := cr.routerGroup.Group(relativePath)

	newCustomRouter := &CustomRouter{
		_router:       cr._router,
		routerGroup:   newGroup,
		log:           cr.log,
		cfg:           cr.cfg,
		storage:       cr.storage,
		botSharedData: cr.botSharedData,
	}

	groupHandler(newCustomRouter)
}

func SetupCustomRouter(
	log *slog.Logger,
	cfg *config.Config,
	storage *sqlite.Storage,
	botSharedData *bot.BotSharedData,
) *CustomRouter {
	router := gin.New()

	return &CustomRouter{
		_router:       router,
		routerGroup:   &router.RouterGroup,
		log:           log,
		cfg:           cfg,
		storage:       storage,
		botSharedData: botSharedData,
	}
}
