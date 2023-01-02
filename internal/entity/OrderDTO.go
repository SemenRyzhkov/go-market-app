package entity

type OrderDTO struct {
	Number string `json:"number"`
	Status string `json:"status"`
	//Accrual    float64     `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
	//UserID     string      `json:"user_id"`
}
