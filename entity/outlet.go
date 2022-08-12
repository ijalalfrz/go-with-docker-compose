package entity

import "time"

type Outlet struct {
	ID         int64     `json:"id"`
	MerchantID int64     `json:"merchant_id"`
	OutletName string    `json:"outlet_name"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
}
