package routes

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type CustomRouter struct {
	_router     *gin.Engine
	routerGroup *gin.RouterGroup
	log         *slog.Logger
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
		_router:     cr._router,
		routerGroup: newGroup,
		log:         cr.log,
	}

	groupHandler(newCustomRouter)
}

func SetupCustomRouter(log *slog.Logger) *CustomRouter {
	router := gin.Default()

	return &CustomRouter{
		_router:     router,
		routerGroup: &router.RouterGroup,
		log:         log,
	}
}
