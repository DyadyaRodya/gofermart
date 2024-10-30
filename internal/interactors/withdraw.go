package interactors

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
	"time"
)

type WithdrawInteractor struct {
	repo              interfaces.Repository
	luhnDomainService interfaces.LuhnService
}

func NewWithdrawInteractor(repo interfaces.Repository, ls interfaces.LuhnService) *WithdrawInteractor {
	return &WithdrawInteractor{
		repo:              repo,
		luhnDomainService: ls,
	}
}

func (i *WithdrawInteractor) Handle(
	ctx context.Context,
	userUUID,
	orderNumber string,
	sum domainmodels.AccrualPoint,
) error {
	if !i.luhnDomainService.Validate(orderNumber) {
		return domainmodels.ErrOrderNumberInvalid
	}
	withdraw := &domainmodels.Withdraw{
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now(),
		UserUUID:    userUUID,
	}

	dbSess, err := i.repo.NewSerializableSession(ctx)
	if err != nil {
		return fmt.Errorf("WithdrawInteractor.repo.NewSerializableSession: %w", err)
	}
	defer dbSess.Close(ctx)

	_, err = dbSess.AddWithdraw(ctx, withdraw)
	if errors.Is(err, domainmodels.ErrNotEnoughPointsToWithdraw) {
		return err
	}
	if errors.Is(err, domainmodels.ErrWithdrawExists) {
		return err
	}
	if err != nil {
		return fmt.Errorf("WithdrawInteractor.dbSess.AddWithdraw: %w", err)
	}

	err = dbSess.Commit(ctx)
	if err != nil {
		return fmt.Errorf("WithdrawInteractor.dbSess.Commit: %w", err)
	}
	return nil
}
