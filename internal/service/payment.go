package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AlmatKozhakhmetov/payment-service/internal/domain"
	"github.com/AlmatKozhakhmetov/payment-service/internal/repository"
)

type PaymentService struct {
	repo repository.IdempotencyRepository
}

func NewPaymentService(repo repository.IdempotencyRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, idempotencyKey string, amount float64, currency string) (*domain.Payment, error) {
	// 1. Проверяем, есть ли готовый ответ
	rec, exists, err := s.repo.Get(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}
	if exists && rec.Status == domain.StatusCompleted {
		return rec.Response, nil
	}

	// 2. Блокируем ключ
	err = s.repo.Lock(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}

	// 3. Создаем платеж
	payment := &domain.Payment{
		ID:        fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		Amount:    amount,
		Currency:  currency,
		Status:    "SUCCESS",
		CreatedAt: time.Now(),
	}

	// 4. Сохраняем результат
	err = s.repo.SaveResponse(ctx, idempotencyKey, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
