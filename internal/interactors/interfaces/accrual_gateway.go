package interfaces

import (
	"context"
	"github.com/DyadyaRodya/gofermart/internal/interactors/dto"
)

type OrderAccrualGateway interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (*dto.OrderAccrual, error)
}
