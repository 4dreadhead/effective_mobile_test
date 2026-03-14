package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"effective_mobile_test/internal/subscriptions/model"

	"gorm.io/gorm"
)

func (u *SubscriptionUsecase) Update(ctx context.Context, id uint64, input UpdateSubscriptionInput) (*model.Subscription, error) {
	sub, err := u.repo.GetSubscription(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if input.ServiceName != nil {
		name := strings.TrimSpace(*input.ServiceName)
		if name == "" {
			return nil, fmt.Errorf("%w: service name is required", ErrValidation)
		}
		sub.ServiceName = name
	}
	if input.MonthlyCost != nil {
		sub.MonthlyCost = *input.MonthlyCost
	}
	if input.From != nil {
		fromDate, err := parseMonthYear(*input.From)
		if err != nil {
			return nil, err
		}
		sub.FromDate = *fromDate
	}
	if input.To != nil {
		toDate, err := parseMonthYear(*input.To)
		if err != nil {
			return nil, err
		}
		sub.ToDate = toDate
	}

	if err = validateUpdate(sub); err != nil {
		return nil, err
	}
	sub, err = u.repo.UpdateSubscription(ctx, *sub)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func validateUpdate(sub *model.Subscription) error {
	if sub.MonthlyCost <= 0 {
		return fmt.Errorf("%w: monthly cost must be positive", ErrValidation)
	}
	if strings.TrimSpace(sub.ServiceName) == "" {
		return fmt.Errorf("%w: service name is required", ErrValidation)
	}
	if sub.FromDate.IsZero() {
		return fmt.Errorf("%w: from date is required", ErrValidation)
	}
	if sub.ToDate != nil && sub.ToDate.Before(sub.FromDate) {
		return fmt.Errorf("%w: to must be on or after from", ErrValidation)
	}

	return nil
}
