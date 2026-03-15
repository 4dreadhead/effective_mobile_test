package presenter

import (
	"fmt"

	"effective_mobile_test/internal/subscriptions/model"
)

type SubscriptionResponse struct {
	ID          uint64  `json:"id"`
	ServiceName string  `json:"service_name"`
	MonthlyCost int     `json:"monthly_cost"`
	UserID      string  `json:"user_id"`
	From        string  `json:"from"`
	To          *string `json:"to,omitempty"`
}

type TotalResponse struct {
	TotalRub int64 `json:"total_rub"`
	HasData  bool  `json:"has_data"`
}

func ToSubscriptionResponse(sub *model.Subscription) SubscriptionResponse {
	from := fmt.Sprintf("%02d.%04d", sub.FromDate.Month(), sub.FromDate.Year())
	var to *string
	if sub.ToDate != nil {
		val := fmt.Sprintf("%02d.%04d", sub.ToDate.Month(), sub.ToDate.Year())
		to = &val
	}
	return SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		MonthlyCost: sub.MonthlyCost,
		UserID:      sub.UserID,
		From:        from,
		To:          to,
	}
}

func ToSubscriptionList(subs []model.Subscription) []SubscriptionResponse {
	out := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		out = append(out, ToSubscriptionResponse(&sub))
	}
	return out
}

func ToTotalResponse(total int64, hasData bool) TotalResponse {
	return TotalResponse{
		TotalRub: total,
		HasData:  hasData,
	}
}
