package model

const (
	InvoiceTransactionType  = "Invoice"
	WithdrawTransactionType = "Withdraw"

	CreatedTransactionStatus = "Created"
)

type Account struct {
	Id            int64   `json:"id"`
	CardNumber    string  `json:"cardNumber"`
	Balance       float64 `json:"balance"`
	FrozenBalance float64 `json:"frozen_balance"`
}

type AccountBalance struct {
	CardNumber    string  `json:"card_number"`
	ActualBalance float64 `json:"actual_balance"`
	FrozenBalance float64 `json:"frozen_balance"`
}
