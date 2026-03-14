package usecase

import (
	"context"
	"fmt"
	"time"
)

func (u *SubscriptionUsecase) TotalCost(ctx context.Context, filter TotalFilter) (int64, bool, error) {
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
		return nil, nil, fmt.Errorf("%w: invalid from interval", ErrValidation)
	}
	toDate, err := parseMonthYear(to)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: invalid to interval", ErrValidation)
	}

	return fromDate, toDate, nil
}
