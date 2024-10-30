package interactors

import (
	"context"
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	interactorsdto "github.com/DyadyaRodya/gofermart/internal/interactors/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
	"time"
)

type RegisterInteractor struct {
	repo            interfaces.Repository
	passwordService interfaces.PasswordService
	loginService    interfaces.LoginService
	uuidGenerator   interfaces.UUIDGenerator
	saltSize        int
}

func NewRegisterInteractor(
	repo interfaces.Repository,
	ps interfaces.PasswordService,
	ls interfaces.LoginService,
	uuidGenerator interfaces.UUIDGenerator,
	saltSize int,
) *RegisterInteractor {
	return &RegisterInteractor{
		repo:            repo,
		passwordService: ps,
		loginService:    ls,
		uuidGenerator:   uuidGenerator,
		saltSize:        saltSize,
	}
}

func (i *RegisterInteractor) Handle(ctx context.Context, creds *interactorsdto.Credentials) (*domainmodels.UserInfo, error) {
	dbSess, err := i.repo.NewSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("RegisterInteractor.repo.NewSession: %w", err)
	}
	defer dbSess.Close(ctx)

	err = i.loginService.Validate(creds.Login)
	if err != nil {
		return nil, err
	}

	userInfo, err := dbSess.GetUserByLogin(ctx, creds.Login)
	if err != nil && !errors.Is(err, domainmodels.ErrUserNotFound) {
		return nil, fmt.Errorf("RegisterInteractor.dbSess.GetUserByLogin: %w", err)
	}

	if userInfo != nil { // same as !errors.Is(err, domainmodels.ErrUserNotFound)
		return nil, domainmodels.ErrLoginTaken
	}

	if !i.passwordService.Validate(creds.Password) {
		return nil, domainmodels.ErrPasswordComplexity
	}

	uuid, err := i.uuidGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("RegisterInteractor.uuidGenerator.Generate: %w", err)
	}

	passwordSalt, err := i.passwordService.GenerateRandomSalt(i.saltSize)
	if err != nil {
		return nil, fmt.Errorf("RegisterInteractor.passwordService.GenerateRandomSalt: %w", err)
	}

	passwordHash, err := i.passwordService.Hash(creds.Password, passwordSalt)
	if err != nil {
		return nil, fmt.Errorf("RegisterInteractor.passwordService.Hash: %w", err)
	}

	userInfo = &domainmodels.UserInfo{
		UUID:         uuid,
		Login:        creds.Login,
		CreatedAt:    time.Now(),
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}

	err = dbSess.AddUser(ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("RegisterInteractor.dbSess.AddUser: %w", err)
	}

	err = dbSess.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("RegisterInteractor.dbSess.Commit: %w", err)
	}
	return userInfo, nil
}
