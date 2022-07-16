package models

type PaymentIntent struct {
	// ID is the unique identifier for the payment intent.
	ID string `json:"id"`
	// Amount is the amount of the payment.
	Amount int64 `json:"amount"`
	// CreatedAt is the date the payment intent was created.
	CreatedAt int64 `json:"created_at"`
	// Captured is true if the payment intent has been captured.
	Captured bool `json:"captured"`
	// Refunded is true if the payment intent has been refunded.
	Refunded bool `json:"refunded"`
	// RefundId is the id of the refund if the payment intent has been refunded.
	RefundId string `json:"refund_id"`
	// ChargeId is the id of the charge if the payment intent has been captured.
	ChargeId string `json:"charge_id"`
}

type PaymentIntentRequest struct {
	// ID is the unique identifier for the payment intent.
	ID string `json:"id"`
	// Amount is the amount of the payment.
	Amount int64 `json:"amount"`
}

type GenericIDResponse struct {
	// ID is the unique identifier for generic response
	ID string `json:"id"`
	// Error is the error message if any
	Error string `json:"error"`
}
