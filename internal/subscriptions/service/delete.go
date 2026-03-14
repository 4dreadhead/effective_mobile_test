package service

import (
	"context"
)

func (u *SubscriptionService) Delete(ctx context.Context, id uint64) error {
	return u.repo.DeleteSubscription(ctx, id)
}
