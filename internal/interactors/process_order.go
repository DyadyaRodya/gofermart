package interactors

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
)

type (
	ProcessOrderErrorLogger func(msg string)
	ProcessOrderInteractor  struct {
		repo    interfaces.Repository
		gateway interfaces.OrderAccrualGateway
	}
)

func NewProcessOrderInteractor(repo interfaces.Repository, gateway interfaces.OrderAccrualGateway) *ProcessOrderInteractor {
	return &ProcessOrderInteractor{
		repo:    repo,
		gateway: gateway,
	}
}

func (i *ProcessOrderInteractor) Run(
	ctx context.Context,
	errorLogger ProcessOrderErrorLogger,
	opChan chan *domainmodels.Order,
) {
	for {
		select {
		case order := <-opChan:
			updatedOrder, err := i.Handle(ctx, order)
			if err != nil {
				errorLogger(fmt.Sprintf("ProcessOrderInteractor.Handle: %v", err))
				go func() { opChan <- order }() // try again later
				continue
			}

			if !updatedOrder.Status.Finished() { // processing yet
				go func() { opChan <- updatedOrder }()
				continue
			}
		case <-ctx.Done():
			return
		}
	}

}

func (i *ProcessOrderInteractor) Handle(
	ctx context.Context,
	order *domainmodels.Order,
) (*domainmodels.Order, error) {
	var orderAccrual *interactorsdto.OrderAccrual = nil
	var err error

	// read gateway only when status unfinished (need only retry db writing after error)
	if !order.Status.Finished() {
		orderAccrual, err = i.gateway.GetOrderAccrual(ctx, order.Number)
		if errors.Is(err, domainmodels.ErrAccrualOrderNotRegistered) {
			// The technical specification doesn't say what to do, we'll try later
			return order, err
		}
		if err != nil {
			return order, fmt.Errorf("ProcessOrderInteractor.gateway.GetOrderAccrual: %v for order %s", err, order.Number)
		}
	}

	if orderAccrual != nil && order.Status == orderAccrual.Status { // no updates
		return order, nil
	}

	if orderAccrual != nil { // have updates
		order.Status = orderAccrual.Status
		order.Accrual = orderAccrual.Accrual
	}

	// update order in db after gateway reading or retry db writing after error
	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return order, fmt.Errorf("ProcessOrderInteractor.repo.NewSession: %v for order %s", err, order.Number)
	}
	defer dbSess.Close(ctx)

	err = dbSess.UpdateOrder(ctx, order)
	if err != nil {
		return order, fmt.Errorf("ProcessOrderInteractor.dbSess.UpdateOrder: %v for order %s", err, order.Number)
	}

	err = dbSess.Commit(ctx)
	if err != nil {
		return order, fmt.Errorf("ProcessOrderInteractor.dbSess.Commit: %v for order %s", err, order.Number)
	}
	return order, nil
}
