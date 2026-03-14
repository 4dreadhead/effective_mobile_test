package repository

import (
	"context"
	apperrors "effective_mobile_test/internal/platform/errors"
	"effective_mobile_test/internal/subscriptions/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateSubscription(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {
	record := &model.Subscription{
		ServiceName: sub.ServiceName,
		UserID:      sub.UserID,
		MonthlyCost: sub.MonthlyCost,
		FromDate:    sub.FromDate,
		ToDate:      sub.ToDate,
	}
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, r.mapError(err)
	}
	return record, nil
}

func (r *PostgresRepository) GetSubscription(ctx context.Context, id uint64) (*model.Subscription, error) {
	var record model.Subscription
	err := r.db.WithContext(ctx).First(&record, id).Error

	if err != nil {
		return nil, r.mapError(err)
	}

	return &record, nil
}

func (r *PostgresRepository) UpdateSubscription(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
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

func (r *PostgresRepository) DeleteSubscription(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Subscription{}, id).Error
}

func (r *PostgresRepository) ListSubscriptions(ctx context.Context, filter ListFilter) ([]model.Subscription, error) {
	query := r.db.WithContext(ctx).Model(&model.Subscription{})
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ServiceName != "" {
		query = query.Where("service_name = ?", filter.ServiceName)
	}

	var rows []model.Subscription
	if err := query.Order("id asc").Find(&rows).Error; err != nil {
		return nil, r.mapError(err)
	}

	out := make([]model.Subscription, 0, len(rows))
	for _, item := range rows {
		out = append(out, item)
	}
	return out, nil
}

func (r *PostgresRepository) TotalCost(ctx context.Context, filter TotalFilter) (int64, bool, error) {
	query := r.db.WithContext(ctx).
		Model(&model.Subscription{}).
		Select("sum(monthly_cost)")

	if filter.ServiceName != "" {
		query = query.Where("service_name = ?", filter.ServiceName)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ToDate != nil {
		query = query.Where("from_date <= ?", filter.ToDate)
	}
	subquery := r.db.Where("to_date IS NULL")
	if filter.FromDate != nil {
		subquery = subquery.Or("to_date >= ?", filter.FromDate)
	}
	query = query.Where(subquery)

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

func (r *PostgresRepository) mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrRecordNotFound
	}
	return fmt.Errorf("repo error: %w", err)
}
