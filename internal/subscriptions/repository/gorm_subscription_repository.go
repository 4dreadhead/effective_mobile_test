package repository

import (
	"context"
	apperrors "effective_mobile_test/internal/platform/errors"
	"effective_mobile_test/internal/subscriptions/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PgSubscriptionRepository struct {
	db *gorm.DB
}

func NewPgSubscriptionRepository(db *gorm.DB) *PgSubscriptionRepository {
	return &PgSubscriptionRepository{db: db}
}

func (r *PgSubscriptionRepository) CreateSubscription(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {
	record := &Subscription{
		ServiceName: sub.ServiceName,
		UserID:      sub.UserID,
		MonthlyCost: sub.MonthlyCost,
		FromDate:    sub.FromDate,
		ToDate:      sub.ToDate,
	}
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, r.mapError(err)
	}

	return record.ToDomainModel(), nil
}

func (r *PgSubscriptionRepository) GetSubscription(ctx context.Context, id uint64) (*model.Subscription, error) {
	var record Subscription
	err := r.db.WithContext(ctx).First(&record, id).Error

	if err != nil {
		return nil, r.mapError(err)
	}

	return record.ToDomainModel(), nil
}

func (r *PgSubscriptionRepository) UpdateSubscription(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	updates := map[string]any{
		"service_name": sub.ServiceName,
		"user_id":      sub.UserID,
		"monthly_cost": sub.MonthlyCost,
		"from_date":    sub.FromDate,
		"to_date":      sub.ToDate,
	}
	err := r.db.WithContext(ctx).Model(sub).
		Where("id = ?", sub.ID).
		Updates(updates).Error

	if err != nil {
		return nil, r.mapError(err)
	}
	return r.GetSubscription(ctx, sub.ID)
}

func (r *PgSubscriptionRepository) DeleteSubscription(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&Subscription{}, id).Error
}

func (r *PgSubscriptionRepository) ListSubscriptions(ctx context.Context, filter ListFilter) ([]model.Subscription, error) {
	query := r.db.WithContext(ctx).Model(&Subscription{})
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ServiceName != "" {
		query = query.Where("service_name = ?", filter.ServiceName)
	}

	var rows []Subscription
	if err := query.Order("id ASC").Find(&rows).Error; err != nil {
		return nil, r.mapError(err)
	}

	out := make([]model.Subscription, 0, len(rows))
	for _, item := range rows {
		out = append(out, *item.ToDomainModel())
	}
	return out, nil
}

func (r *PgSubscriptionRepository) TotalCost(ctx context.Context, filter TotalFilter) (int64, bool, error) {
	query := r.db.WithContext(ctx).
		Model(&Subscription{}).
		Select("SUM(monthly_cost)")

	if filter.ServiceName != "" {
		query = query.Where("service_name = ?", filter.ServiceName)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ToDate != nil {
		query = query.Where("from_date <= ?", filter.ToDate)
	}
	if filter.FromDate != nil {
		query = query.Where(r.db.Where("to_date IS NULL").Or("to_date >= ?", filter.FromDate))
	}

	var total *int64
	err := query.Scan(&total).Error
	if err != nil {
		return 0, false, r.mapError(err)
	}
	if total == nil {
		return 0, false, nil
	}

	return *total, true, nil
}

func (r *PgSubscriptionRepository) mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrRecordNotFound
	}
	return fmt.Errorf("repo error: %w", err)
}
