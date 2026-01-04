package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/shared/coreErrors"
)

type subscriptionPlanStorage struct {
	db Database
}

func (u *subscriptionPlanStorage) MustInit() {
	query := `CREATE TABLE IF NOT EXISTS subscription_plans(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		price INTEGER NOT NULL,
		duration_days INTEGER NOT NULL,
		vpn_server_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vpn_server_id) REFERENCES vpn_servers(id)
	);`

	_, err := u.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init subscription_plans table: %w", err))
	}
}

// GetAllSubscriptionPlans returns all available subscription plans
func (sp *subscriptionPlanStorage) GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans = make([]models.SubscriptionPlan, 0)

	query := "SELECT * FROM subscription_plans"
	err := sp.db.Select(&plans, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription plans: %w", err)
	}

	return plans, nil
}

// GetSubscriptionPlanByID gets a subscription plan by ID
func (sp *subscriptionPlanStorage) GetSubscriptionPlanByID(subscriptionPlanId int64) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan

	query := "SELECT * FROM subscription_plans WHERE id = ?"
	err := sp.db.Get(&plan, query, subscriptionPlanId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, coreErrors.ErrorNotFound
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return &plan, nil
}

// CreateSubscriptionPlan adds a new subscription plan to the database
func (sp *subscriptionPlanStorage) CreateSubscriptionPlan(plan *models.CreateSubscriptionPlan) (*models.SubscriptionPlan, error) {
	query := `
    INSERT INTO subscription_plans (name, description, price, duration_days, vpn_server_id)
    VALUES (:name, :description, :price, :duration_days, :vpn_server_id) RETURNING *
    `
	rows, err := sp.db.NamedQuery(query, plan)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: subscription_plans.name") {
			return nil, coreErrors.ErrorAlreadyExist
		}

		return nil, fmt.Errorf("failed to create subscription plan: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("rows next error")
	}

	var createdSubscriptionPlan models.SubscriptionPlan

	err = rows.StructScan(&createdSubscriptionPlan)

	if err != nil {
		return nil, fmt.Errorf("failed to scan of created subscription plan: %w", err)
	}

	return &createdSubscriptionPlan, nil
}

// UpdateSubscriptionPlan updates an existing subscription plan in the database
func (sp *subscriptionPlanStorage) UpdateSubscriptionPlan(updatePlanId int64, updatePlanPayload *models.UpdateSubscriptionPlan) (*models.SubscriptionPlan, error) {
	// TODO: for change vpn server id need to change clients configs

	query := fmt.Sprintf(`
    UPDATE subscription_plans
    SET name = :name,
		description = :description,
		price = :price,
		duration_days = :duration_days
    WHERE id = %d RETURNING *
    `, updatePlanId)

	rows, err := sp.db.NamedQuery(query, updatePlanPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription plan: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, coreErrors.ErrorNotFound
	}

	var updatedPlan models.SubscriptionPlan

	err = rows.StructScan(&updatedPlan)

	if err != nil {
		return nil, fmt.Errorf("struct scan error: %w", err)
	}

	return &updatedPlan, nil
}

// DeleteSubscriptionPlanByID deletes the subscription plan by ID
func (sp *subscriptionPlanStorage) DeleteSubscriptionPlanByID(id int64) error {
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
		return coreErrors.ErrorNotFound
	}

	return nil
}
