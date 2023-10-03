package model

const (
	StatusCreated = "Created"
	StatusError   = "Error"
	StatusSuccess = "Success"
)

type Transaction struct {
	CardNumber string  `json:"card_number"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
}
