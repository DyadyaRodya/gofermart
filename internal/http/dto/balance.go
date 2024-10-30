package dto

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
)

type BalanceDTO struct {
	Current   domainmodels.AccrualPoint `json:"current"`
	Withdrawn domainmodels.AccrualPoint `json:"withdrawn"`
}

func FromBalanceInfo(info *interactorsdto.BalanceInfo) *BalanceDTO {
	return &BalanceDTO{
		Current:   info.Current,
		Withdrawn: info.Withdrawn,
	}
}
