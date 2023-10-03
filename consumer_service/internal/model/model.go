package model

const (
	StatusCreated = "Created"
	StatusError   = "Error"
	StatusSuccess = "Success"
)

type Transaction struct {
	Card     string
	Currency string
	Amount   float64
	Status   string
}
