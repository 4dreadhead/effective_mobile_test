package usecase

import (
	"context"
	"effective_mobile_test/internal/subscriptions/model"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub model.Subscription) (*model.Subscription, error)
	GetSubscription(ctx context.Context, id uint64) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	DeleteSubscription(ctx context.Context, id uint64) error
	ListSubscriptions(ctx context.Context, filter ListFilter) ([]model.Subscription, error)
	TotalCost(ctx context.Context, filter TotalFilter) (int64, bool, error)
}

type CreateSubscriptionInput struct {
	ServiceName string
	MonthlyCost int
	UserID      string
	From        string
	To          string
}

type UpdateSubscriptionInput struct {
	ServiceName *string
	MonthlyCost *int
	From        *string
	To          *string
}

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

type SubscriptionUsecase struct {
	repo SubscriptionRepository
}

func NewSubscriptionUsecase(repo SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: repo}
}

func parseMonthYear(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month")
	}
	if year < 1 {
		return nil, fmt.Errorf("invalid year")
	}

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return &date, nil
}
