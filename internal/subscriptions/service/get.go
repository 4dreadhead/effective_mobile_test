package service

import (
	"context"
	"effective_mobile_test/internal/subscriptions/model"
)

func (u *SubscriptionService) Get(ctx context.Context, id uint64) (*model.Subscription, error) {
	return u.repo.GetSubscription(ctx, id)
}
