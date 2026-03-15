package orderservice

import "errors"

var (
	ErrOrderNotFound               = errors.New("order not found")
	ErrCodeExpired                 = errors.New("code expired")
	ErrCodeNotFound                = errors.New("code not found")
	ErrAmountOfUsage               = errors.New("amount of usage is zero")
	ErrAmountOfBooks               = errors.New("insufficient number of books in stock")
	ErrEmptyCart                   = errors.New("empty cart")
	ErrAlreadyPaid                 = errors.New("order already paid")
	ErrFailedToPay                 = errors.New("failed to pay order")
	ErrDeliveryDestinationRequired = errors.New("either address or shop_id must be provided")
	ErrDeliveryDestinationConflict = errors.New("address and shop_id cannot be set together")
)
