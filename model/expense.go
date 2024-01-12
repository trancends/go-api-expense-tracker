package model

import "time"

type transactionType string

const (
	CREDIT transactionType = "CREDIT"
	DEBIT  transactionType = "DEBIT"
)

type Expense struct {
	ID              string          `json:"id"`
	Date            string          `json:"date"`
	Amount          float64         `json:"amount"`
	TransactionType transactionType `json:"transactionType"`
	Balance         float64         `json:"balance"`
	Description     string          `json:"description"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
}
