package services

import (
	"bizarre-vpn-api/internal/storage/models"
	"bizarre-vpn-api/internal/storage/repositories"
	"bizarre-vpn-api/pkg/custom_errors"
	"errors"
	"fmt"
)

// CreatePlan creates a new subscription plan
func CreatePlan(plan *models.SubscriptionPlan) (int64, error) {
	if plan.DurationMonths <= 0 {
		return 0, fmt.Errorf("duration months must be greater than zero")
	}
	if plan.Price < 0 {
		return 0, fmt.Errorf("price must be non-negative")
	}
	if plan.DeviceLimit < 0 {
		return 0, fmt.Errorf("device limit must be greater or equal to zero")
	}

	return repositories.CreateSubscriptionPlan(plan)
}

// GetPlan gets a subscription plan by ID
func GetPlan(id int64) (*models.SubscriptionPlan, error) {
	plan, err := repositories.GetSubscriptionPlanByID(id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPlanNotFound) {
			return nil, custom_errors.ErrPlanNotFound
		}
		return nil, fmt.Errorf("error getting subscription plan by id: %v", err)
	}

	if plan == nil {
		return nil, fmt.Errorf("subscription plan not found")
	}

	return plan, nil
}

// GetAllPlans returns all available subscription plans
func GetAllPlans() ([]models.SubscriptionPlan, error) {
	plan, err := repositories.GetAllSubscriptionPlans()
	if err != nil {
		return nil, fmt.Errorf("error getting all subscription plans: %v", err)
	}

	return plan, nil
}

// UpdatePlan updates an existing subscription plan
func UpdatePlan(plan *models.SubscriptionPlan) (*models.SubscriptionPlan, error) {
	if plan.DurationMonths <= 0 {
		return nil, fmt.Errorf("duration months must be greater than zero")
	}
	if plan.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}
	if plan.DeviceLimit < 0 {
		return nil, fmt.Errorf("device limit must be greater or equal to zero")
	}

	err := repositories.UpdateSubscriptionPlan(plan)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPlanNotFound) {
			return nil, custom_errors.ErrPlanNotFound
		}
		return nil, fmt.Errorf("error updating subscription plan: %v", err)
	}

	return plan, nil
}

// DeletePlan deletes the subscription plan by ID
func DeletePlan(id int64) error {
	err := repositories.DeleteSubscriptionPlanByID(id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPlanNotFound) {
			return custom_errors.ErrPlanNotFound
		}
		return fmt.Errorf("error deleting subscription plan by id: %v", err)
	}

	return nil
}
