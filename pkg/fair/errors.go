package fair

import "errors"

var (
	ErrNotFound          = errors.New("Street Fair Not Found")
	ErrInvalidStreetFair = errors.New("Invalid Street Fair")
	ErrInternal          = errors.New("InternalServerError")
)
