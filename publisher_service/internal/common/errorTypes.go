package errorTypes

import "errors"

var (
	ErrJSONUnmarshalling    = errors.New("error while trying to unmarshal json")
	ErrMarshallingToJSON    = errors.New("error while trying to marshal to json")
	ErrNegativeRefillAmount = errors.New("error: refill amount cannot be negative")
	ErrNatsPublishing       = errors.New("error while trying to publish data to the queue")
)
