package usecase

import (
	"context"
	"errors"

	"effective_mobile_test/internal/subscriptions/model"

	"gorm.io/gorm"
)

func (u *SubscriptionUsecase) Get(ctx context.Context, id uint64) (*model.Subscription, error) {
	sub, err := u.repo.GetSubscription(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return sub, err
}
