package interfaces

import "bizarre-vpn-api/internal/models"

type SubscriptionPlanStorage interface {
	GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error)
	GetSubscriptionPlanByID(id int64) (*models.SubscriptionPlan, error)
	CreateSubscriptionPlan(plan *models.CreateSubscriptionPlan) (*models.SubscriptionPlan, error)
	UpdateSubscriptionPlan(updatePlanId int64, updatePlan *models.UpdateSubscriptionPlan) (*models.SubscriptionPlan, error)
	DeleteSubscriptionPlanByID(id int64) error
}
