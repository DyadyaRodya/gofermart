package interactors

import (
	"context"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
)

type GetWithdrawalsInteractor struct {
	repo interfaces.Repository
}

func NewGetWithdrawalsInteractor(repo interfaces.Repository) *GetWithdrawalsInteractor {
	return &GetWithdrawalsInteractor{
		repo: repo,
	}
}

func (i *GetWithdrawalsInteractor) Handle(ctx context.Context, userUUID string) ([]*domainmodels.Withdraw, error) {
	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetWithdrawalsInteractor.repo.NewSession: %w", err)
	}
	defer dbSess.Close(ctx)

	withdrawals, err := dbSess.GetWithdrawalsByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("GetWithdrawalsInteractor.dbSess.GetWithdrawalsByUserUUID: %w", err)
	}
	return withdrawals, nil
}
