package dto

const (
	InvoiceTransactionType  = "Invoice"
	WithdrawTransactionType = "Withdraw"

	CreatedTransactionStatus = "Created"
	ErrorTransactionStatus   = "Error"
	SuccessTransactionStatus = "Success"
)

type TransactionDTO struct {
	CardNumber string  `json:"card_number"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
}
