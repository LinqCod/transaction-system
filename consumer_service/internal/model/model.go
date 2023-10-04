package model

const (
	WithdrawTransactionType = "Withdraw"
	InvoiceTransactionType  = "Invoice"

	StatusError   = "Error"
	StatusSuccess = "Success"
)

type Account struct {
	Id            int64   `json:"id"`
	CardNumber    string  `json:"cardNumber"`
	Balance       float64 `json:"balance"`
	FrozenBalance float64 `json:"frozen_balance"`
}

type Transaction struct {
	CardNumber string  `json:"card_number"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
}
