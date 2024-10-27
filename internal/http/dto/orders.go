package dto

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"time"
)

type OrderDTO struct {
	Number     string                     `json:"number"`
	Status     string                     `json:"status"`
	Accrual    *domainmodels.AccrualPoint `json:"accrual,omitempty"`
	UploadedAt string                     `json:"uploaded_at"`
}

func FromOrderInfo(order *domainmodels.Order) *OrderDTO {
	dto := &OrderDTO{
		Number:     order.Number,
		Status:     string(order.Status),
		UploadedAt: order.UploadedAt.Format(time.RFC3339),
	}
	if order.Accrual > 0 {
		accrual := order.Accrual
		dto.Accrual = &accrual
	}
	return dto
}
