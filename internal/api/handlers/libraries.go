package handlers

import (
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /libraries/backend-types [get]
func (h *LibrariesHandler) GetBackendTypesList(c *gin.Context) {
	const op = "handlers.libraries.GetBackendTypesList"

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

// GetProtocolsList returns protocols list
// @Summary Get protocols list
// @Description Get protocols list like LibraryItems
// @Security token
// @scope.admin only administrative information
// @Tags Libraries
// @Produce json
// @Success 200 {object} []models.LibraryItem "LibraryItems"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /libraries/protocols [get]
func (h *LibrariesHandler) GetProtocolsList(c *gin.Context) {
	const op = "handlers.libraries.GetProtocolsList"

	log := h.Log.With(
		slog.String("op", op),
	)

	list, err := h.LibrariesService.GetProtocolsList()

	if err != nil {
		log.Error("getting list error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetProtocolsListByBackendTypeId returns protocols list by
// @Summary Get protocols list by backend Type Id
// @Description Get protocols list like LibraryItems
// @Security token
// @scope.admin only administrative information
// @Tags Libraries
// @Param backendTypeId query string true "Backend Type Id"
// @Produce json
// @Success 200 {object} []models.LibraryItem "LibraryItems"
// @Failure 404 {object} ErrorResponse "BackendTypeId is not exist"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /libraries/protocols-by-backend-type-id [get]
func (h *LibrariesHandler) GetProtocolsListByBackendTypeId(c *gin.Context) {
	const op = "handlers.libraries.GetProtocolsListByBackendTypeId"

	log := h.Log.With(
		slog.String("op", op),
	)

	backendTypeIdStr, isOk := c.GetQuery("backendTypeId")

	if !isOk {
		log.Info("query parameter 'backendTypeId' is not provided")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "query parameter 'backendTypeId' is required"})
		c.Abort()
		return
	}

	backendTypeId, err := strconv.ParseInt(backendTypeIdStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("convert error backendTypeIdStr: %w", err)
		log.Info(err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	list, err := h.LibrariesService.GetProtocolsListByBackendTypeId(backendTypeId)

	if err != nil {
		if isNotExist := strings.Contains(err.Error(), "is not exist"); isNotExist {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		log.Error("getting list error", sl.Err(err))

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, list)
}
