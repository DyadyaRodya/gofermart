package accrual

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/DyadyaRodya/gofermart/pkg/accrualsdk"
)

func (g *Gateway) fromAccrualInfo(info *accrualsdk.OrderAccrualInfo) *interactorsdto.OrderAccrual {
	var status domainmodels.OrderStatus
	switch info.Status {
	case accrualsdk.OrderAccrualInfoStatusRegistered:
		status = domainmodels.OrderStatusNew
	case accrualsdk.OrderAccrualInfoStatusProcessing:
		status = domainmodels.OrderStatusProcessing
	case accrualsdk.OrderAccrualInfoStatusProcessed:
		status = domainmodels.OrderStatusProcessed
	case accrualsdk.OrderAccrualInfoStatusInvalid:
		status = domainmodels.OrderStatusInvalid
	}

	var accrualPoints domainmodels.AccrualPoint = 0
	if info.Accrual != nil {
		accrualPoints = domainmodels.AccrualPointFromFloat64(*info.Accrual)
	}

	return &interactorsdto.OrderAccrual{
		OrderNumber: info.OrderNumber,
		Status:      status,
		Accrual:     accrualPoints,
	}
}
