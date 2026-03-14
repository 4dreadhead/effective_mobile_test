package presenter

import (
	"fmt"

	"effective_mobile_test/internal/subscriptions/model"
	"effective_mobile_test/internal/subscriptions/view"
)

func ToSubscriptionResponse(sub *model.Subscription) view.SubscriptionResponse {
	from := fmt.Sprintf("%02d.%04d", sub.FromDate.Month(), sub.FromDate.Year())
	var to *string
	if sub.ToDate != nil {
		val := fmt.Sprintf("%02d.%04d", sub.ToDate.Month(), sub.ToDate.Year())
		to = &val
	}
	return view.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		MonthlyCost: sub.MonthlyCost,
		UserID:      sub.UserID,
		From:        from,
		To:          to,
	}
}

func ToSubscriptionList(subs []model.Subscription) []view.SubscriptionResponse {
	out := make([]view.SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		out = append(out, ToSubscriptionResponse(&sub))
	}
	return out
}

func ToTotalResponse(total int64, hasData bool) view.TotalResponse {
	return view.TotalResponse{
		TotalRub: total,
		HasData:  hasData,
	}
}
