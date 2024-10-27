package services

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"slices"
)

type (
	LoginConfig struct {
		MinLen       int
		MaxLen       int
		AllowedChars []rune
	}
	LoginDomainService struct {
		config *LoginConfig
	}
)

func NewLoginDomainService(config *LoginConfig) *LoginDomainService {
	return &LoginDomainService{
		config: config,
	}
}

func (l *LoginDomainService) Validate(login string) error {
	if len(login) < l.config.MinLen {
		return domainmodels.ErrLoginTooShort
	}
	if len(login) > l.config.MaxLen {
		return domainmodels.ErrLoginTooLong
	}

	for _, char := range login {
		if !slices.Contains(l.config.AllowedChars, char) {
			return domainmodels.ErrLoginChars
		}
	}

	return nil
}
