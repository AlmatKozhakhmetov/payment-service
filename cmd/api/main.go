package main

import (
	"log"
	"net/http"

	"github.com/AlmatKozhakhmetov/payment-service/internal/handler"
	"github.com/AlmatKozhakhmetov/payment-service/internal/repository"
	"github.com/AlmatKozhakhmetov/payment-service/internal/service"
)

func main() {
	// Инициализируем слои Clean Architecture
	repo := repository.NewMemoryIdempotencyRepository()
	paymentService := service.NewPaymentService(repo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Регистрируем HTTP-эндпоинт
	http.HandleFunc("/api/v1/payments", paymentHandler.ProcessPayment)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
