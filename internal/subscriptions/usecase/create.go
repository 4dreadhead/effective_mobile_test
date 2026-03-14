package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"effective_mobile_test/internal/subscriptions/model"
)

func (u *SubscriptionUsecase) Create(ctx context.Context, input CreateSubscriptionInput) (*model.Subscription, error) {
	fromDate, toDate, err := validateCreate(input)
	if err != nil {
		return nil, err
	}

	serviceName := strings.TrimSpace(input.ServiceName)
	sub := model.Subscription{
		ServiceName: serviceName,
		UserID:      input.UserID,
		MonthlyCost: input.MonthlyCost,
		FromDate:    *fromDate,
		ToDate:      toDate,
	}
	return u.repo.CreateSubscription(ctx, sub)
}

func validateCreate(input CreateSubscriptionInput) (*time.Time, *time.Time, error) {
	if strings.TrimSpace(input.ServiceName) == "" {
		return nil, nil, fmt.Errorf("%w: service name is required", ErrValidation)
	}
	if input.MonthlyCost <= 0 {
		return nil, nil, fmt.Errorf("%w: monthly cost must be positive", ErrValidation)
	}

	fromDate, err := parseMonthYear(input.From)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: invalid from period", ErrValidation)
	}

	var toDate *time.Time

	if input.To != "" {
		toDate, err = parseMonthYear(input.To)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: invalid to period", ErrValidation)
		}
		if toDate.Before(*fromDate) {
			return nil, nil, fmt.Errorf("%w: to must be on or after from", ErrValidation)
		}
	}

	return fromDate, toDate, nil
}
