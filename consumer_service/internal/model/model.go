package model

const (
	InvoiceTransactionType  = "Invoice"
	WithdrawTransactionType = "Withdraw"

	StatusCreated = "Created"
	StatusError   = "Error"
	StatusSuccess = "Success"
)

type Account struct {
	Id         int64   `json:"id"`
	CardNumber string  `json:"cardNumber"`
	Balance    float64 `json:"balance"`
}

type Transaction struct {
	CardNumber string  `json:"card_number"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
}
