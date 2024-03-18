package domain

import "time"

type Account struct {
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	Balance       string    `json:"balance"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
