package model

import "time"

type GroupOmzetPerMerchant struct {
	MerchantID          *int64    `json:"merchant_id"`
	MerchantName        string    `json:"merchant_name"`
	NumberOftransaction int       `json:"number_of_transaction"`
	Omzet               int64     `json:"omzet"`
	TransactionDate     string    `json:"-"`
	TransactionMonth    string    `json:"-"`
	TransactionYear     string    `json:"-"`
	TransactionFullDate time.Time `json:"transaction_full_date"`
}

type GroupOmzetPerOutlet struct {
	OutletID            *int64    `json:"outlet_id"`
	OutletName          string    `json:"outlet_name"`
	MerchantName        string    `json:"merchant_name"`
	NumberOftransaction int       `json:"number_of_transaction"`
	Omzet               int64     `json:"omzet"`
	TransactionDate     string    `json:"-"`
	TransactionMonth    string    `json:"-"`
	TransactionYear     string    `json:"-"`
	TransactionFullDate time.Time `json:"transaction_full_date"`
}
