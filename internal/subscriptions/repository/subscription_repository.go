package repository

import (
	"context"
	"effective_mobile_test/internal/subscriptions/model"
	"time"
)

type ListFilter struct {
	UserID      string
	ServiceName string
}

type TotalFilter struct {
	UserID      string
	ServiceName string
	From        string
	To          string
	FromDate    *time.Time
	ToDate      *time.Time
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub model.Subscription) (*model.Subscription, error)
	GetSubscription(ctx context.Context, id uint64) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	DeleteSubscription(ctx context.Context, id uint64) error
	ListSubscriptions(ctx context.Context, filter ListFilter) ([]model.Subscription, error)
	TotalCost(ctx context.Context, filter TotalFilter) (int64, bool, error)
}
