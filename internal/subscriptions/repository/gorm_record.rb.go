package repository

import (
	"effective_mobile_test/internal/subscriptions/model"
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ServiceName string
	UserID      string
	MonthlyCost int
	FromDate    time.Time
	ToDate      *time.Time
}

func (rec *Subscription) ToDomainModel() *model.Subscription {
	return &model.Subscription{
		ID:          uint64(rec.ID),
		ServiceName: rec.ServiceName,
		UserID:      rec.UserID,
		MonthlyCost: rec.MonthlyCost,
		FromDate:    rec.FromDate,
		ToDate:      rec.ToDate,
	}
}
