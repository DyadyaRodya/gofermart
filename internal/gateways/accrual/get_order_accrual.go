package accrual

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/DyadyaRodya/gofermart/pkg/accrualsdk"
)

func (g *Gateway) GetOrderAccrual(ctx context.Context, orderNumber string) (*interactorsdto.OrderAccrual, error) {
	for {
		orderAccrualInfo, err := g.client.GetOrderAccrual(ctx, orderNumber)

		if errors.Is(err, context.Canceled) {
			g.errorLogger("Gateway.GetOrderAccrual context canceled")
			return nil, err
		}

		if errors.Is(err, accrualsdk.ErrAccrualClientRetryLater) {
			g.errorLogger(fmt.Sprintf("Gateway.GetOrderAccrual retrying request for order %s", orderNumber))
			continue
		}

		if err != nil {
			g.errorLogger(fmt.Sprintf(
				"Gateway.GetOrderAccrual unexpected request error %v for order %s", err, orderNumber,
			))
			return nil, errors.Join(err, ErrAccrualProcessorGateway)
		}

		if orderAccrualInfo == nil {
			return nil, domainmodels.ErrAccrualOrderNotRegistered
		}

		orderAccrual := g.fromAccrualInfo(orderAccrualInfo)
		return orderAccrual, nil
	}

}
