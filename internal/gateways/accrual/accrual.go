package accrual

import (
	"errors"
	"fmt"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/pkg/accrualsdk"
	"net/url"
)

var ErrAccrualProcessorGateway = errors.Join(domainmodels.ErrInternalServer, errors.New("accrual processor gateway error"))

type GatewayErrorLogger func(msg string)

type Gateway struct {
	config      *Config
	gatewayURL  *url.URL
	client      *accrualsdk.AccrualClient
	errorLogger GatewayErrorLogger
}

func NewGateway(config *Config) (*Gateway, error) {
	gatewayURL, err := url.Parse(config.Host)
	if err != nil {
		return nil, fmt.Errorf("unable to parse gateway url: %w", err)
	}
	return &Gateway{
		config:     config,
		gatewayURL: gatewayURL,
		client:     accrualsdk.NewAccrualClient(gatewayURL, config.SkipTLS),
	}, nil
}
