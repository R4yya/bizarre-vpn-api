package services

import (
	"errors"
	"fmt"
	"log/slog"
	"unicode/utf8"

	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/services/interfaces"
	"bizarre-vpn-api/internal/shared/coreErrors"
)

type SubscriptionPlanService struct {
	log     *slog.Logger
	storage interfaces.SubscriptionPlanStorage
}

func NewSubscriptionPlanService(
	log *slog.Logger,
	subscriptionPlanStorage interfaces.SubscriptionPlanStorage,
) *SubscriptionPlanService {
	return &SubscriptionPlanService{
		storage: subscriptionPlanStorage,
		log:     log,
	}
}

// GetAllPlans returns all available subscription plans
func (s *SubscriptionPlanService) GetAllPlans() ([]models.SubscriptionPlan, error) {
	op := "internal.services.subscriptionPlanService.GetAllPlans"

	log := s.log.With(slog.String("op", op))

	plans, err := s.storage.GetAllSubscriptionPlans()
	if err != nil {
		return nil, fmt.Errorf("error getting all subscription plans: %w", err)
	}

	log.Debug("plans list", slog.Any("plans", plans))

	return plans, nil
}

// GetPlan gets a subscription plan by ID
func (s *SubscriptionPlanService) GetPlan(id int64) (*models.SubscriptionPlan, error) {
	plan, err := s.storage.GetSubscriptionPlanByID(id)
	if err != nil {
		if errors.Is(err, coreErrors.ErrorNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("error getting subscription plan by id: %w", err)
	}

	if plan == nil {
		return nil, fmt.Errorf("subscription plan not found")
	}

	return plan, nil
}

// CreatePlan creates a new subscription plan
func (s *SubscriptionPlanService) CreatePlan(createPlan *models.CreateSubscriptionPlan) (*models.SubscriptionPlan, error) {
	if nameLength := utf8.RuneCountInString(createPlan.Name); nameLength < 5 || nameLength >= 250 {
		return nil, coreErrors.ValidationError{
			Msg:    "must be longer than 5 and shorted than 250",
			Entity: "subscription plan",
			Fields: []string{
				"name",
			},
		}
	}

	if createPlan.DurationDays <= 0 {
		return nil, coreErrors.ValidationError{
			Msg:    "must be greater than zero",
			Entity: "subscription plan",
			Fields: []string{
				"duration days",
			},
		}
	}

	if createPlan.Price < 0 {
		return nil, coreErrors.ValidationError{
			Msg:    "must be non-negative",
			Entity: "subscription plan",
			Fields: []string{
				"price",
			},
		}
	}

	if createPlan.VpnServerID == 0 {
		return nil, coreErrors.ValidationError{
			Msg:    "field is required",
			Entity: "subscription plan",
			Fields: []string{
				"vpnServerId",
			},
		}
	}

	return s.storage.CreateSubscriptionPlan(createPlan)
}

// UpdatePlan updates an existing subscription plan
func (s *SubscriptionPlanService) UpdatePlan(updatePlanId int64, updatePlan *models.UpdateSubscriptionPlan) (*models.SubscriptionPlan, error) {
	if nameLength := utf8.RuneCountInString(updatePlan.Name); nameLength < 5 || nameLength >= 250 {
		return nil, coreErrors.ValidationError{
			Msg:    "must be longer than 5 and shorted than 250",
			Entity: "subscription plan",
			Fields: []string{
				"name",
			},
		}
	}

	if updatePlan.DurationDays <= 0 {
		return nil, coreErrors.ValidationError{
			Msg:    "must be greater than zero",
			Entity: "subscription plan",
			Fields: []string{
				"duration days",
			},
		}
	}

	if updatePlan.Price < 0 {
		return nil, coreErrors.ValidationError{
			Msg:    "must be non-negative",
			Entity: "subscription plan",
			Fields: []string{
				"price",
			},
		}
	}

	updatedSubscriptionPlan, err := s.storage.UpdateSubscriptionPlan(updatePlanId, updatePlan)
	if err != nil {
		if errors.Is(err, coreErrors.ErrorNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("error updating subscription plan: %w", err)
	}

	return updatedSubscriptionPlan, nil
}

// DeletePlan deletes the subscription plan by ID
func (s *SubscriptionPlanService) DeletePlan(id int64) error {
	err := s.storage.DeleteSubscriptionPlanByID(id)
	if err != nil {
		if errors.Is(err, coreErrors.ErrorNotFound) {
			return err
		}
		return fmt.Errorf("error deleting subscription plan by id: %w", err)
	}

	return nil
}
