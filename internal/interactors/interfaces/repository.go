package interfaces

import (
	"context"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
)

type RepoSession interface {
	GetUserByLogin(ctx context.Context, login string) (*domainmodels.UserInfo, error)
	AddUser(ctx context.Context, user *domainmodels.UserInfo) error

	AddOrder(ctx context.Context, order *domainmodels.Order) (*domainmodels.Order, error)
	UpdateOrder(ctx context.Context, order *domainmodels.Order) error
	GetOrdersByUserUUID(ctx context.Context, userUUID string) ([]*domainmodels.Order, error)

	AddWithdraw(ctx context.Context, withdraw *domainmodels.Withdraw) (*domainmodels.Withdraw, error)
	GetWithdrawalsByUserUUID(ctx context.Context, userUUID string) ([]*domainmodels.Withdraw, error)

	GetTransactionSumInfoByUserUUID(ctx context.Context, userUUID string) (*interactorsdto.TransactionSumInfo, error)

	Commit(ctx context.Context) error
	Close(ctx context.Context) error
}

type Repository interface {
	NewSession(ctx context.Context) (RepoSession, error)
	NewSerializableSession(ctx context.Context) (RepoSession, error)
}
