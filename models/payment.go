package models

import "time"

type Payment struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	UserID      int       `json:"user_id"`
}
