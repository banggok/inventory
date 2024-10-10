package usecase

import "errors"

// ErrProductNotFound is returned when a product is not found in the repository
var ErrProductNotFound = errors.New("product not found")
