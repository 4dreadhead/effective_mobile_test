package model

import (
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
