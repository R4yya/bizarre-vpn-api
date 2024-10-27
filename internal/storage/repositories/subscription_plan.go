package repositories

import (
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
	"fmt"
)

// CreateSubscriptionPlan adds a new subscription plan to the database
func CreateSubscriptionPlan(plan *models.SubscriptionPlan) (int64, error) {
	query := `
    INSERT INTO subscription_plans (country, name, description, duration_months, data_limit_gb, speed_limit_mbps, device_limit, price)
    VALUES (:country, :name, :description, :duration_months, :data_limit_gb, :speed_limit_mbps, :device_limit, :price)
    `
	result, err := storage.GetDB().NamedExec(query, plan)
	if err != nil {
		return 0, fmt.Errorf("failed to create subscription plan: %w", err)
	}

	planID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return planID, nil
}

// GetSubscriptionPlanByID gets a subscription plan by ID
func GetSubscriptionPlanByID(id int64) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan

	query := "SELECT * FROM subscription_plans WHERE id = ?"
	err := storage.GetDB().Get(&plan, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return &plan, nil
}

// GetAllSubscriptionPlans returns all available subscription plans
func GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan

	query := "SELECT * FROM subscription_plans"
	err := storage.GetDB().Select(&plans, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription plans: %w", err)
	}

	return plans, nil
}

// UpdateSubscriptionPlan updates an existing subscription plan in the database
func UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error {
	query := `
    UPDATE subscription_plans
    SET country = :country, name = :name, description = :description,
        duration_months = :duration_months, data_limit_gb = :data_limit_gb,
        speed_limit_mbps = :speed_limit_mbps, device_limit = :device_limit,
        price = :price
    WHERE id = :id
    `
	_, err := storage.GetDB().NamedExec(query, plan)
	if err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}
	return nil
}

// DeleteSubscriptionPlanByID deletes the subscription plan by ID
func DeleteSubscriptionPlanByID(id int64) error {
	query := "DELETE FROM subscription_plans WHERE id = ?"
	result, err := storage.GetDB().Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription plan found with ID %d", id)
	}

	return nil
}
