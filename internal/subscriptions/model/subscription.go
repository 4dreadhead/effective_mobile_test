package model

import (
	"time"
)

type Subscription struct {
	ID          uint64
	ServiceName string
	UserID      string
	MonthlyCost int
	FromDate    time.Time
	ToDate      *time.Time
}
