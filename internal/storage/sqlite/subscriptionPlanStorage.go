package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type SubscriptionPlanStorage struct {
	db Database
}

func (u *SubscriptionPlanStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS subscription_plans(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
		-- language_code TEXT NOT NULL
		-- created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		-- role VARCHAR(255),
		-- refresh_token TEXT(500)
	);`

	//CREATE INDEX IF NOT EXISTS idx_id ON users(id)

	_, err := u.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init subscription_plans table: %v", err))
	}
}

// GetAllSubscriptionPlans returns all available subscription plans
func (sp *SubscriptionPlanStorage) GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan

	query := "SELECT * FROM subscription_plans"
	err := sp.db.Select(&plans, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription plans: %w", err)
	}

	return plans, nil
}

// GetSubscriptionPlanByID gets a subscription plan by ID
func (sp *SubscriptionPlanStorage) GetSubscriptionPlanByID(id int64) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan

	query := "SELECT * FROM subscription_plans WHERE id = ?"
	err := sp.db.Get(&plan, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrPlanNotFound
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return &plan, nil
}

// CreateSubscriptionPlan adds a new subscription plan to the database
func (sp *SubscriptionPlanStorage) CreateSubscriptionPlan(plan *models.SubscriptionPlan) (int64, error) {
	query := `
    INSERT INTO subscription_plans (country, name, description, duration_months, data_limit_gb, speed_limit_mbps, device_limit, price)
    VALUES (:country, :name, :description, :duration_months, :data_limit_gb, :speed_limit_mbps, :device_limit, :price)
    `
	result, err := sp.db.NamedExec(query, plan)
	if err != nil {
		return 0, fmt.Errorf("failed to create subscription plan: %w", err)
	}

	planID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return planID, nil
}

// UpdateSubscriptionPlan updates an existing subscription plan in the database
func (sp *SubscriptionPlanStorage) UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error {
	query := `
    UPDATE subscription_plans
    SET country = :country, name = :name, description = :description,
        duration_months = :duration_months, data_limit_gb = :data_limit_gb,
        speed_limit_mbps = :speed_limit_mbps, device_limit = :device_limit,
        price = :price
    WHERE id = :id
    `

	result, err := sp.db.NamedExec(query, plan)
	if err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return storage.ErrPlanNotFound
	}
	return nil
}

// DeleteSubscriptionPlanByID deletes the subscription plan by ID
func (sp *SubscriptionPlanStorage) DeleteSubscriptionPlanByID(id int64) error {
	query := "DELETE FROM subscription_plans WHERE id = ?"
	result, err := sp.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return storage.ErrPlanNotFound
	}

	return nil
}
