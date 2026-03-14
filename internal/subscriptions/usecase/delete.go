package usecase

import (
	"context"
)

func (u *SubscriptionUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.DeleteSubscription(ctx, id)
}
