package entity

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
