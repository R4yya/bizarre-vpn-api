package handlers

import (
	"bizarre-vpn-api/internal/storage/models"
	"bizarre-vpn-api/internal/storage/services"
	"bizarre-vpn-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SubscriptionPlanRequest struct {
	Country        string  `json:"country" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	DurationMonths int     `json:"duration_months" binding:"required"`
	DataLimitGB    *int    `json:"dataLimitGB"`
	SpeedLimitMbps *int    `json:"speedLimitMbps"`
	DeviceLimit    int     `json:"deviceLimit" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
}

// CreatePlanHandler creates a new subscription plan
func CreatePlanHandler(c *gin.Context) {
	var req SubscriptionPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}

	plan := &models.SubscriptionPlan{
		Country:        req.Country,
		Name:           req.Name,
		Description:    req.Description,
		DurationMonths: req.DurationMonths,
		DataLimitGB:    req.DataLimitGB,
		SpeedLimitMbps: req.SpeedLimitMbps,
		DeviceLimit:    req.DeviceLimit,
		Price:          req.Price,
	}

	planID, err := services.CreatePlan(plan)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	plan.ID = planID
	c.JSON(http.StatusCreated, plan)
}

// GetPlanHandler returns the subscription plan by ID
func GetPlanHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid plan ID"})
		return
	}

	plan, err := services.GetPlan(id)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusNotFound, MessageResponse{Message: "plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// GetAllPlansHandler returns all available subscription plans
func GetAllPlansHandler(c *gin.Context) {
	plans, err := services.GetAllPlans()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// UpdatePlanHandler updates the subscription plan by ID
func UpdatePlanHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid plan ID"})
		return
	}

	var req SubscriptionPlanRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}

	plan := &models.SubscriptionPlan{
		ID:             id,
		Country:        req.Country,
		Name:           req.Name,
		Description:    req.Description,
		DurationMonths: req.DurationMonths,
		DataLimitGB:    req.DataLimitGB,
		SpeedLimitMbps: req.SpeedLimitMbps,
		DeviceLimit:    req.DeviceLimit,
		Price:          req.Price,
	}

	updatedPlan, err := services.UpdatePlan(plan)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPlan)
}

// DeletePlanHandler deletes the subscription plan by ID
func DeletePlanHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid plan ID"})
		return
	}

	if err := services.DeletePlan(id); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "plan successfully deleted"})
}
