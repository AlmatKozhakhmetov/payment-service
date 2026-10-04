package domain

import (
	"errors"
	"time"
)

var ErrDuplicateRequest = errors.New("duplicate request processing")

type IdempotencyStatus string

const (
	StatusProcessing IdempotencyStatus = "PROCESSING"
	StatusCompleted  IdempotencyStatus = "COMPLETED"
)

type Payment struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createde_at"`
}

type IdempotencyRecord struct {
	Key       string            `json:"key"`
	Status    IdempotencyStatus `json:"status"`
	Response  *Payment          `json:"paymentresponse,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}
