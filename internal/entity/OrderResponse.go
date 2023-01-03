package entity

type OrderResponse struct {
	Order   int     `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
