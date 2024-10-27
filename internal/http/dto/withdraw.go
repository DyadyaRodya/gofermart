package dto

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"time"
)

type WithdrawDTO struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

func (w *WithdrawDTO) ConvertSum() domainmodels.AccrualPoint {
	return domainmodels.AccrualPointFromFloat64(w.Sum)
}

type WithdrawalDTO struct {
	OrderNumber string                    `json:"order"`
	Sum         domainmodels.AccrualPoint `json:"sum"`
	ProcessedAt string                    `json:"processedAt"`
}

func FromWithdrawalInfo(withdrawal *domainmodels.Withdraw) *WithdrawalDTO {
	return &WithdrawalDTO{
		OrderNumber: withdrawal.OrderNumber,
		Sum:         withdrawal.Sum,
		ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
	}
}
