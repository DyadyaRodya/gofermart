package interactors

import (
	"context"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
)

type GetOrdersInteractor struct {
	repo interfaces.Repository
}

func NewGetOrdersInteractor(repo interfaces.Repository) *GetOrdersInteractor {
	return &GetOrdersInteractor{
		repo: repo,
	}
}

func (i *GetOrdersInteractor) Handle(ctx context.Context, userUUID string) ([]*domainmodels.Order, error) {
	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetOrdersInteractor.repo.NewSession: %w", err)
	}
	defer dbSess.Close(ctx)

	orders, err := dbSess.GetOrdersByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("GetOrdersInteractor.dbSess.GetOrdersByUserUUID: %w", err)
	}
	return orders, nil
}
