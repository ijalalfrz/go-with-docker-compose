package entity

import "time"

type Merchant struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	MerchantName string    `json:"merchant_name"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
}
