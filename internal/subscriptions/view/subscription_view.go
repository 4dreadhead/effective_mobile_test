package view

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
