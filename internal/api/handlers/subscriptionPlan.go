package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"bizarre-vpn-api/internal/shared/logger/sl"
)

type SubscriptionPlanHandler struct {
	Log                     *slog.Logger
	SubscriptionPlanService *services.SubscriptionPlanService
}

// GetAllPlansHandler returns all available subscription plans
// @Summary Get all subscription plans
// @Description Retrieves all available subscription plans
// @Tags SubscriptionPlans
// @Produce json
// @Success 200 {array} models.SubscriptionPlan "List of all plans"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscription-plans [get]
func (h *SubscriptionPlanHandler) GetAllPlansHandler(c *gin.Context) {
	const op = "subscription_plan.GetAllPlansHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	plans, err := h.SubscriptionPlanService.GetAllPlans()
	if err != nil {
		log.Error("getting plans error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetPlanHandler returns the subscription plan by ID
// @Summary Get a plan subscription by ID
// @Description Retrieves a subscription plan by its unique identifier
// @Tags SubscriptionPlans
// @Produce json
// @Param id path int true "Plan ID"
// @Success 200 {object} models.SubscriptionPlan "Plan details"
// @Failure 400 {object} ErrorResponse "Invalid plan ID"
// @Failure 404 {object} ErrorResponse "Plan not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscription-plans/{id} [get]
func (h *SubscriptionPlanHandler) GetPlanHandler(c *gin.Context) {
	const op = "subscription_plan.GetPlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		log.Error("parsing uri param id error", sl.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid plan ID"})
		return
	}

	plan, err := h.SubscriptionPlanService.GetPlan(id)
	if errors.Is(err, coreErrors.ErrorNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		log.Error("getting plan error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// CreatePlanHandler creates a new subscription plan
// @Summary Create a subscription plan
// @Description Creates a new subscription plan with the specified details
// @Tags SubscriptionPlans
// @Accept json
// @Produce json
// @Param plan body models.CreateSubscriptionPlan true "Plan Data"
// @Success 201 {object} models.SubscriptionPlan "Successfully created plan"
// @Failure 400 {object} ErrorResponse "Invalid request or missing required parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscription-plans [post]
func (h *SubscriptionPlanHandler) CreatePlanHandler(c *gin.Context) {
	const op = "handlers.subscription_plan.CreatePlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	var req models.CreateSubscriptionPlan

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("parsing request json error", sl.Err(err))

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdPlan, err := h.SubscriptionPlanService.CreatePlan(&req)
	if err != nil {
		var validation coreErrors.ValidationError

		if errors.As(err, &validation) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, coreErrors.ErrorAlreadyExist) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		log.Error("creating plan error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, createdPlan)
}

// UpdatePlanHandler updates the subscription plan by ID
// @Summary Update a subscription plan
// @Description Updates the subscription plan with the specified details
// @Tags SubscriptionPlans
// @Accept json
// @Produce json
// @Param id path int true "Plan ID"
// @Param plan body models.UpdateSubscriptionPlan true "Updated Plan Data"
// @Success 200 {object} models.SubscriptionPlan "Successfully updated plan"
// @Failure 400 {object} ErrorResponse "Invalid plan ID or request body"
// @Failure 404 {object} ErrorResponse "Plan not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscription-plans/{id} [put]
func (h *SubscriptionPlanHandler) UpdatePlanHandler(c *gin.Context) {
	const op = "subscription_plan.UpdatePlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	idStr := c.Param("id")
	updatePlanId, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		log.Error("parsing uri param id error", sl.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid plan ID"})
		return
	}

	var req models.UpdateSubscriptionPlan
	if err = c.ShouldBindJSON(&req); err != nil {
		log.Error("parsing request json error", sl.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedPlan, err := h.SubscriptionPlanService.UpdatePlan(updatePlanId, &req)
	if err != nil {
		var validation coreErrors.ValidationError

		if errors.As(err, &validation) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, coreErrors.ErrorNotFound) {
			h.Log.Info("update plan not found", sl.Err(err))
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		h.Log.Error("update plan error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, updatedPlan)
}

// DeletePlanHandler deletes the subscription plan by ID
// @Summary Delete a subscription plan by ID
// @Description Deletes the subscription plan with the specified ID
// @Tags SubscriptionPlans
// @Accept json
// @Produce json
// @Param id path int true "Plan ID"
// @Success 200 {object} MessageResponse "Successfully deleted plan"
// @Failure 400 {object} ErrorResponse "Invalid plan ID"
// @Failure 404 {object} ErrorResponse "Plan not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscription-plans/{id} [delete]
func (h *SubscriptionPlanHandler) DeletePlanHandler(c *gin.Context) {
	const op = "subscription_plan.UpdatePlanHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("parsing uri param id error", sl.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid plan ID"})
		return
	}

	if err = h.SubscriptionPlanService.DeletePlan(id); err != nil {
		if errors.Is(err, coreErrors.ErrorNotFound) {
			h.Log.Info("delete plan not found err")
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		h.Log.Error("delete plan error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "plan successfully deleted"})
}
