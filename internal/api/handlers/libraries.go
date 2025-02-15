package handlers

import (
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LibrariesHandler struct {
	Log              *slog.Logger
	LibrariesService *services.LibrariesService
}

// GetBackendTypesList returns backendTypes list
// @Summary Get backendTypes list
// @Description Get backendTypes list like LibraryItems
// @Security token
// @scope.admin only administrative information
// @Tags Libraries
// @Produce json
// @Success 200 {object} []models.LibraryItem "LibraryItems"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /libraries/backend-types [get]
func (h *LibrariesHandler) GetBackendTypesList(c *gin.Context) {
	const op = "handlers.auth.AuthorizeWithInitData"

	log := h.Log.With(
		slog.String("op", op),
	)

	list, err := h.LibrariesService.GetBackendTypesList()

	if err != nil {
		log.Error("getting list error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, list)
}
