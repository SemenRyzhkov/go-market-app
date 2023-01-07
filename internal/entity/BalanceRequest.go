package entity

type BalanceRequest struct {
	Current   float32 `json:"current,omitempty"`
	Withdrawn float32 `json:"withdrawn,omitempty"`
}
