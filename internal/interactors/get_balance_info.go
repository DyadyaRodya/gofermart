package interactors

import (
	"context"
	"fmt"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
)

type GetBalanceInteractor struct {
	repo interfaces.Repository
}

func NewGetBalanceInteractor(repo interfaces.Repository) *GetBalanceInteractor {
	return &GetBalanceInteractor{
		repo: repo,
	}
}

func (i *GetBalanceInteractor) Handle(ctx context.Context, userUUID string) (*interactorsdto.BalanceInfo, error) {
	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetBalanceInteractor.repo.NewSession: %w", err)
	}
	defer dbSess.Close(ctx)

	transactionSumInfo, err := dbSess.GetTransactionSumInfoByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("GetBalanceInteractor.dbSess.GetTransactionSumInfoByUserUUID: %w", err)
	}
	return transactionSumInfo.BalanceInfo(), nil
}
