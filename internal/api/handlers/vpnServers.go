package handlers

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VpnServersHandler struct {
	Log               *slog.Logger
	VpnServersService *services.VpnServersService
}

type CreateBodyInput struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	AdapterHost      string `json:"adapterHost"`
	AdapterPort      uint16 `json:"adapterPort"`
	Country          string `json:"country"`
	LnkBackendTypeId int64  `json:"lnkBackendTypeId"`
}

// GetExpandedList returns VpnServerExpandedItem list
// @Summary Get VpnServerExpandedItem list
// @Description Get VpnServerExpandedItem list
// @Security token
// @scope.admin only administrative information
// @Tags VpnServers
// @Produce json
// @Success 200 {object} []models.VpnServerExpandedItem "VpnServerExpandedItems"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /vpn-servers [get]
func (h *VpnServersHandler) GetExpandedList(c *gin.Context) {
	op := "internal.handlers.vpnServersHandler.GetExpandedList"

	log := h.Log.With(slog.String("op", op))

	list, err := h.VpnServersService.GetExpandedList()

	if err != nil {
		log.Error("error of handling get list", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, list)
	c.Abort()
}

// GetExpandedItem returns VpnServerExpandedItem by id
// @Summary Get VpnServerItem by id
// @Description Get VpnServerItem by id
// @Security token
// @scope.admin only administrative information
// @Tags VpnServers
// @Param id path int true "VpnServerItem Id"
// @Produce json
// @Success 200 {object} models.VpnServerExpandedItem "VpnServerExpandedItem"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "VpnServerItem is not exist"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /vpn-servers/{id} [get]
func (h *VpnServersHandler) GetExpandedItem(c *gin.Context) {
	op := "internal.handlers.vpnServersHandler.GetExpandedItem"

	log := h.Log.With(slog.String("op", op))

	vpnServerIdStr := c.Param("id")

	vpnServerId, err := strconv.ParseInt(vpnServerIdStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("convert error vpnServerIdStr: %w", err)
		log.Info(err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	item, err := h.VpnServersService.GetExpandedItemById(vpnServerId)

	if err != nil {
		if isNotFound := errors.Is(err, coreErrors.ErrorNotFound); isNotFound {
			log.Info(err.Error())
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "vpn server not found"})
			c.Abort()
			return
		}

		log.Error("error of handling get list", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem returns created VpnServerItem
// @Summary Create VpnServerItem
// @Description Create VpnServerItem by body data
// @Security token
// @scope.admin only administrative information
// @Tags VpnServers
// @Accept json
// @Produce json
// @Param VpnServerInfo body CreateBodyInput true "VpnServerInfo Data"
// @Success 201 {object} models.VpnServerExpandedItem "Successful Created VpnServerExpandedItem"
// @Failure 400 {object} ErrorResponse "Not Unique name"
// @Failure 400 {object} ErrorResponse "Not Unique host and port"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /vpn-servers [post]
func (h *VpnServersHandler) CreateItem(c *gin.Context) {
	op := "internal.handlers.vpnServersHandler.CreateItem"

	log := h.Log.With(slog.String("op", op))

	var body CreateBodyInput

	err := c.ShouldBindJSON(&body)

	if err != nil {
		log.Error("parsing request json error", sl.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	log.Debug("body", slog.Any("CreateBodyInput", body))

	vpnServer := models.VpnServerItem{
		Name:             body.Name,
		Description:      body.Description,
		AdapterHost:      body.AdapterHost,
		AdapterPort:      body.AdapterPort,
		LnkBackendTypeId: body.LnkBackendTypeId,
		Country:          body.Country,
	}

	createdVpnServerId, err := h.VpnServersService.CreateItem(&vpnServer)

	if err != nil {
		log.Error("creating new vpnServer", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	vpnServerExpanded, err := h.VpnServersService.GetExpandedItemById(createdVpnServerId)

	if err != nil {
		log.Error("error of getting created vpnServerExpanded", sl.Err(err))

		var validation coreErrors.ValidationError

		if errors.As(err, &validation) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, vpnServerExpanded)
}

// Delete Item By ID
// @Summary Delete VpnServerItem by id
// @Description Delete VpnServerItem by id
// @Security token
// @scope.admin only administrative information
// @Tags VpnServers
// @Param id path int true "VpnServerItem Id"
// @Produce json
// @Success 200 {object} MessageResponse "Success"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "VpnServerItem is not exist"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /vpn-servers/{id} [delete]
func (h *VpnServersHandler) DeleteItem(c *gin.Context) {
	op := "internal.handlers.vpnServersHandler.DeleteItem"

	log := h.Log.With(slog.String("op", op))

	vpnServerIdStr := c.Param("id")

	vpnServerId, err := strconv.ParseInt(vpnServerIdStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("convert error vpnServerIdStr: %w", err)
		log.Info(err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	err = h.VpnServersService.DeleteItem(vpnServerId)

	if err != nil {
		if isNotFound := errors.Is(err, coreErrors.ErrorNotFound); isNotFound {
			log.Info(err.Error())
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "vpn server not found"})
			c.Abort()
			return
		}

		log.Error("delete by id error", sl.Err(err))

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{Message: "VpnServer successful deleted"})
}
