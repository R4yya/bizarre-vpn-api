package services

import (
	"errors"
	"fmt"

	intStorage "bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type SubscriptionPlanStorage interface {
	GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id int64) (*models.SubscriptionPlan, error)
	CreateSubscriptionPlan(plan *models.SubscriptionPlan) (int64, error)
	UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error
	DeleteSubscriptionPlanByID(id int64) error
}

type SubscriptionPlanService struct {
	storage SubscriptionPlanStorage
}

// GetAllPlans returns all available subscription plans
func (s *SubscriptionPlanService) GetAllPlans() ([]models.SubscriptionPlan, error) {
	plan, err := s.storage.GetAllSubscriptionPlans()
	if err != nil {
		return nil, fmt.Errorf("error getting all subscription plans: %w", err)
	}

	return plan, nil
}

// GetPlan gets a subscription plan by ID
func (s *SubscriptionPlanService) GetPlan(id int64) (*models.SubscriptionPlan, error) {
	plan, err := s.storage.GetSubscriptionPlanByID(id)
	if err != nil {
		if errors.Is(err, intStorage.ErrPlanNotFound) {
			return nil, intStorage.ErrPlanNotFound
		}
		return nil, fmt.Errorf("error getting subscription plan by id: %w", err)
	}

	if plan == nil {
		return nil, fmt.Errorf("subscription plan not found")
	}

	return plan, nil
}

// CreatePlan creates a new subscription plan
func (s *SubscriptionPlanService) CreatePlan(plan *models.SubscriptionPlan) (int64, error) {
	if plan.DurationMonths <= 0 {
		return 0, fmt.Errorf("duration months must be greater than zero")
	}
	if plan.Price < 0 {
		return 0, fmt.Errorf("price must be non-negative")
	}
	if plan.DeviceLimit < 0 {
		return 0, fmt.Errorf("device limit must be greater or equal to zero")
	}

	return s.storage.CreateSubscriptionPlan(plan)
}

// UpdatePlan updates an existing subscription plan
func (s *SubscriptionPlanService) UpdatePlan(plan *models.SubscriptionPlan) (*models.SubscriptionPlan, error) {
	if plan.DurationMonths <= 0 {
		return nil, fmt.Errorf("duration months must be greater than zero")
	}
	if plan.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}
	if plan.DeviceLimit < 0 {
		return nil, fmt.Errorf("device limit must be greater or equal to zero")
	}

	err := s.storage.UpdateSubscriptionPlan(plan)
	if err != nil {
		if errors.Is(err, intStorage.ErrPlanNotFound) {
			return nil, intStorage.ErrPlanNotFound
		}
		return nil, fmt.Errorf("error updating subscription plan: %w", err)
	}

	return plan, nil
}

// DeletePlan deletes the subscription plan by ID
func (s *SubscriptionPlanService) DeletePlan(id int64) error {
	err := s.storage.DeleteSubscriptionPlanByID(id)
	if err != nil {
		if errors.Is(err, intStorage.ErrPlanNotFound) {
			return intStorage.ErrPlanNotFound
		}
		return fmt.Errorf("error deleting subscription plan by id: %w", err)
	}

	return nil
}
