package usecase

import (
	"context"

	"effective_mobile_test/internal/subscriptions/model"
)

func (u *SubscriptionUsecase) List(ctx context.Context, filter ListFilter) ([]model.Subscription, error) {
	return u.repo.ListSubscriptions(ctx, filter)
}
