package interactors

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
	"time"
)

type AddOrderInteractor struct {
	repo               interfaces.Repository
	luhnDomainService  interfaces.LuhnService
	orderProcessorChan chan *domainmodels.Order
}

func NewAddOrderInteractor(repo interfaces.Repository, ls interfaces.LuhnService, opChan chan *domainmodels.Order) *AddOrderInteractor {
	return &AddOrderInteractor{
		repo:               repo,
		luhnDomainService:  ls,
		orderProcessorChan: opChan,
	}
}

func (i *AddOrderInteractor) Handle(ctx context.Context, userUUID, orderNumber string) error {
	if !i.luhnDomainService.Validate(orderNumber) {
		return domainmodels.ErrOrderNumberInvalid
	}
	order := &domainmodels.Order{
		Number:     orderNumber,
		Status:     domainmodels.OrderStatusNew,
		UploadedAt: time.Now(),
		UserUUID:   userUUID,
	}

	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return fmt.Errorf("AddOrderInteractor.repo.NewSession: %w", err)
	}
	defer dbSess.Close(ctx)

	order, err = dbSess.AddOrder(ctx, order)
	if errors.Is(err, domainmodels.ErrOrderExists) {
		return err
	}
	if err != nil {
		return fmt.Errorf("AddOrderInteractor.dbSess.AddOrder: %w", err)
	}
	go func() { i.orderProcessorChan <- order }() // counting accrual in another handler,
	// TODO handler should read unprocessed on start from db in case of service restarting

	err = dbSess.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AddOrderInteractor.dbSess.Commit: %w", err)
	}
	return nil
}
