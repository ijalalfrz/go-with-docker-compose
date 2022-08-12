package entity

import "time"

type Transaction struct {
	User      *User     `json:"user"`
	Merchant  *Merchant `json:"merchant"`
	Outlet    *Outlet   `json:"outlet"`
	ID        int64     `json:"id"`
	OutletID  int64     `json:"outlet_id"`
	BillTotal float64   `json:"bill_total"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}
