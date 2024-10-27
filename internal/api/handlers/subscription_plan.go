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
	DurationMonths int     `json:"durationMonths"`
	DataLimitGB    *int    `json:"dataLimitGb"`
	SpeedLimitMbps *int    `json:"speedLimitMbps"`
	DeviceLimit    int     `json:"deviceLimit"`
	Price          float64 `json:"price"`
}

// CreatePlanHandler creates a new subscription plan
// @Summary Create a subscription plan
// @Description Creates a new subscription plan with the specified details
// @Tags Subscription Plans
// @Accept json
// @Produce json
// @Param plan body SubscriptionPlanRequest true "Plan Data"
// @Success 201 {object} models.SubscriptionPlan "Successfully created plan"
// @Failure 400 {object} MessageResponse "Invalid request or missing required parameters"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /plans [post]
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
// @Summary Get a plan subscription by ID
// @Description Retrieves a subscription plan by its unique identifier
// @Tags Subscription Plans
// @Accept json
// @Produce json
// @Param id path int true "Plan ID"
// @Success 200 {object} models.SubscriptionPlan "Plan details"
// @Failure 400 {object} MessageResponse "Invalid plan ID"
// @Failure 404 {object} MessageResponse "Plan not found"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /plans/{id} [get]
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
// @Summary Get all subscription plans
// @Description Retrieves all available subscription plans
// @Tags Subscription Plans
// @Accept json
// @Produce json
// @Success 200 {array} models.SubscriptionPlan "List of all plans"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /plans [get]
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
// @Summary Update a subscription plan
// @Description Updates the subscription plan with the specified details
// @Tags Subscription Plans
// @Accept json
// @Produce json
// @Param id path int true "Plan ID"
// @Param plan body SubscriptionPlanRequest true "Updated Plan Data"
// @Success 200 {object} models.SubscriptionPlan "Successfully updated plan"
// @Failure 400 {object} MessageResponse "Invalid plan ID or request body"
// @Failure 404 {object} MessageResponse "Plan not found"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /plans/{id} [put]
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
// @Summary Delete a subscription plan by ID
// @Description Deletes the subscription plan with the specified ID
// @Tags Subscription Plans
// @Accept json
// @Produce json
// @Param id path int true "Plan ID"
// @Success 200 {object} MessageResponse "Successfully deleted plan"
// @Failure 400 {object} MessageResponse "Invalid plan ID"
// @Failure 404 {object} MessageResponse "Plan not found"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /plans/{id} [delete]
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
