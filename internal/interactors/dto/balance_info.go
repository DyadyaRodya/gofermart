package dto

import domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"

type TransactionSumInfo struct {
	OrderAccruals domainmodels.AccrualPoint
	Withdrawn     domainmodels.AccrualPoint
}

type BalanceInfo struct {
	Current   domainmodels.AccrualPoint
	Withdrawn domainmodels.AccrualPoint
}

func (t TransactionSumInfo) BalanceInfo() *BalanceInfo {
	return &BalanceInfo{
		Current:   t.OrderAccruals - t.Withdrawn,
		Withdrawn: t.Withdrawn,
	}
}
