package service

import (
	"effective_mobile_test/internal/subscriptions/repository"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func parseMonthYear(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month")
	}
	if year < 1 {
		return nil, fmt.Errorf("invalid year")
	}

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return &date, nil
}
