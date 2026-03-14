package model

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ID          uint64
	ServiceName string
	UserID      string
	MonthlyCost int
	FromDate    time.Time
	ToDate      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
