package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type SubscriptionPlanHandler struct {
	Log                     *slog.Logger
	subscriptionPlanService *services.SubscriptionPlanService
}

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

// GetAllPlansHandler returns all available subscription plans
// @Summary Get all subscription plans
// @Description Retrieves all available subscription plans
// @Tags Subscription Plans
// @Accept json
// @Produce json
// @Success 200 {array} models.SubscriptionPlan "List of all plans"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /plans [get]
func (h *SubscriptionPlanHandler) GetAllPlansHandler(c *gin.Context) {
	const op = "subscription_plan.GetAllPlansHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	plans, err := h.subscriptionPlanService.GetAllPlans()
	if err != nil {
		log.Error("getting plans error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
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
func (h *SubscriptionPlanHandler) GetPlanHandler(c *gin.Context) {
	const op = "subscription_plan.GetPlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		log.Error("parsing uri param id error", sl.Err(err))
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid plan ID"})
		return
	}

	plan, err := h.subscriptionPlanService.GetPlan(id)
	if errors.Is(err, storage.ErrPlanNotFound) {
		c.JSON(http.StatusNotFound, MessageResponse{Message: err.Error()})
		return
	} else if err != nil {
		log.Error("getting plan error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
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
func (h *SubscriptionPlanHandler) CreatePlanHandler(c *gin.Context) {
	const op = "handlers.subscription_plan.CreatePlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	var req SubscriptionPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("parsing request json error", sl.Err(err))

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

	planID, err := h.subscriptionPlanService.CreatePlan(plan)
	if err != nil {
		log.Error("creating plan error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	plan.ID = planID
	c.JSON(http.StatusCreated, plan)
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
func (h *SubscriptionPlanHandler) UpdatePlanHandler(c *gin.Context) {
	const op = "subscription_plan.UpdatePlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("parsing uri param id error", sl.Err(err))
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid plan ID"})
		return
	}

	var req SubscriptionPlanRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		log.Error("parsing request json error", sl.Err(err))
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

	updatedPlan, err := h.subscriptionPlanService.UpdatePlan(plan)
	if err != nil {
		if errors.Is(err, storage.ErrPlanNotFound) {
			h.Log.Info("update plan not found", sl.Err(err))
			c.JSON(http.StatusNotFound, MessageResponse{Message: err.Error()})
			return
		}

		h.Log.Error("update plan error", sl.Err(err))
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
func (h *SubscriptionPlanHandler) DeletePlanHandler(c *gin.Context) {
	const op = "subscription_plan.UpdatePlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("parsing uri param id error", sl.Err(err))
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "invalid plan ID"})
		return
	}

	if err = h.subscriptionPlanService.DeletePlan(id); err != nil {
		if errors.Is(err, storage.ErrPlanNotFound) {
			h.Log.Info("delete plan not found err")
			c.JSON(http.StatusNotFound, MessageResponse{Message: err.Error()})
			return
		}
		h.Log.Error("delete plaln error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "plan successfully deleted"})
}
