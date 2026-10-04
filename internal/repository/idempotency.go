package repository

import (
	"context"
	"sync"

	"github.com/AlmatKozhakhmetov/payment-service/internal/domain"
)

type IdempotencyRepository interface {
	Get(ctx context.Context, key string) (*domain.IdempotencyRecord, bool, error)
	Lock(ctx context.Context, key string) error
	SaveResponse(ctx context.Context, key string, payment *domain.Payment) error
}

type MemoryIdempotencyRepository struct {
	mu      sync.RWMutex
	storage map[string]*domain.IdempotencyRecord
}

func NewMemoryIdempotencyRepository() *MemoryIdempotencyRepository {
	return &MemoryIdempotencyRepository{
		storage: make(map[string]*domain.IdempotencyRecord),
	}
}

func (r *MemoryIdempotencyRepository) Get(ctx context.Context, key string) (*domain.IdempotencyRecord, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, exists := r.storage[key]
	return rec, exists, nil
}

func (r *MemoryIdempotencyRepository) Lock(ctx context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rec, exists := r.storage[key]; exists {
		if rec.Status == domain.StatusProcessing {
			return domain.ErrDuplicateRequest
		}
		return nil
	}

	r.storage[key] = &domain.IdempotencyRecord{
		Key:    key,
		Status: domain.StatusProcessing,
	}
	return nil
}

func (r *MemoryIdempotencyRepository) SaveResponse(ctx context.Context, key string, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rec, exists := r.storage[key]; exists {
		rec.Status = domain.StatusCompleted
		rec.Response = payment
	}
	return nil
}
