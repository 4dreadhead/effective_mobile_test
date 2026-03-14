package service

import (
	"context"
	"effective_mobile_test/internal/subscriptions/repository"

	"effective_mobile_test/internal/subscriptions/model"
)

func (u *SubscriptionService) List(ctx context.Context, filter repository.ListFilter) ([]model.Subscription, error) {
	return u.repo.ListSubscriptions(ctx, filter)
}
