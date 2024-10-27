package interactors

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
)

type LoginInteractor struct {
	repo            interfaces.Repository
	passwordService interfaces.PasswordService
}

func NewLoginInteractor(repo interfaces.Repository, ps interfaces.PasswordService, saltSize int) *LoginInteractor {
	return &LoginInteractor{
		repo:            repo,
		passwordService: ps,
	}
}

func (i *LoginInteractor) Handle(
	ctx context.Context,
	creds *interactorsdto.Credentials,
) (*domainmodels.UserInfo, error) {
	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("LoginInteractor.repo.NewSession: %w", err)
	}
	defer dbSess.Close(ctx)

	userInfo, err := dbSess.GetUserByLogin(ctx, creds.Login)
	if err != nil && !errors.Is(err, domainmodels.ErrUserNotFound) {
		return nil, fmt.Errorf("LoginInteractor.dbSess.GetUserByLogin: %w", err)
	}

	if userInfo == nil || err != nil { // errors.Is(err, domainmodels.ErrUserNotFound)
		return nil, domainmodels.ErrWrongCredentials
	}

	passwordMatch, err := i.passwordService.Compare(creds.Password, userInfo.PasswordHash, userInfo.PasswordSalt)
	if err != nil {
		return nil, fmt.Errorf("LoginInteractor.passwordService.Compare: %w", err)
	}
	if !passwordMatch {
		return nil, domainmodels.ErrWrongCredentials
	}

	return userInfo, nil
}
