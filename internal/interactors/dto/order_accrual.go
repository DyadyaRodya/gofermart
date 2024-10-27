package dto

import domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"

type OrderAccrual struct {
	OrderNumber string
	Status      domainmodels.OrderStatus
	Accrual     domainmodels.AccrualPoint
}
