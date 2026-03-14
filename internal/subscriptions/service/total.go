package service

import (
	"context"
	apperrors "effective_mobile_test/internal/platform/errors"
	"effective_mobile_test/internal/subscriptions/repository"
	"fmt"
	"time"
)

func (u *SubscriptionService) TotalCost(ctx context.Context, filter repository.TotalFilter) (int64, bool, error) {
	fromDate, toDate, err := validateFilterPeriod(filter.From, filter.To)
	if err != nil {
		return 0, false, err
	}
	filter.FromDate = fromDate
	filter.ToDate = toDate
	return u.repo.TotalCost(ctx, filter)
}

func validateFilterPeriod(from, to string) (*time.Time, *time.Time, error) {
	fromDate, err := parseMonthYear(from)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: invalid from interval", apperrors.ErrInvalidFields)
	}
	toDate, err := parseMonthYear(to)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: invalid to interval", apperrors.ErrInvalidFields)
	}

	return fromDate, toDate, nil
}
