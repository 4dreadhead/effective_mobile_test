package usecase

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func (u *SubscriptionUsecase) Delete(ctx context.Context, id uint64) error {
	if err := u.repo.DeleteSubscription(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
