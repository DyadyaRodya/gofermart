package config

import (
	"flag"
	domainservices "github.com/DyadyaRodya/gofermart/internal/domain/services"
	accrualgateway "github.com/DyadyaRodya/gofermart/internal/gateways/accrual"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress        string
	AccrualServerAddress string
	LogLevel             string
	DSN                  string
	SaltSize             int

	LoginConfig              *domainservices.LoginConfig
	PasswordComplexityConfig *domainservices.PasswordComplexityConfig

	AccrualGatewayConfig *accrualgateway.Config
}

func InitConfigFromCMD(defaultServerAddress, defaultAccrualServerAddress, defaultLogLevel string, defaultSaltSize int) *Config {
	serverAddress := flag.String("a", defaultServerAddress, "server address to bind")
	accrualServerAddress := flag.String("r", defaultAccrualServerAddress, "accrual server address")
	logLevel := flag.String("l", defaultLogLevel, "log level")
	dsn := flag.String("d", "", "database connection string")
	flag.Parse()

	if envServerAddress := os.Getenv("RUN_ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	if envAccrualServerAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualServerAddress != "" {
		accrualServerAddress = &envAccrualServerAddress
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel = &envLogLevel
	}
	if envDSN := os.Getenv("DATABASE_URI"); envDSN != "" {
		dsn = &envDSN
	}
	saltSize := defaultSaltSize
	if envSaltSize := os.Getenv("SALT_SIZE"); envSaltSize != "" {
		saltSize, _ = strconv.Atoi(envSaltSize)
	}

	return &Config{
		ServerAddress:        *serverAddress,
		AccrualServerAddress: *accrualServerAddress,
		LogLevel:             *logLevel,
		DSN:                  *dsn,
		SaltSize:             saltSize,
		LoginConfig: &domainservices.LoginConfig{
			MinLen: 4, // for e2e tests
			MaxLen: 50,
			AllowedChars: []rune{
				'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
				'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
				'0', '1', '2', '3', '4', '5', '7', '8', '9', '_',
			},
		},
		PasswordComplexityConfig: &domainservices.PasswordComplexityConfig{
			Length:          0, // for e2e tests
			NumberOfDigits:  0, // for e2e tests
			NumberOfUpper:   0, // for e2e tests
			NumberOfLower:   0, // for e2e tests
			NumberOfSpecial: 0, // for e2e tests
		},
		AccrualGatewayConfig: &accrualgateway.Config{
			Host:    *accrualServerAddress,
			SkipTLS: true,
		},
	}
}
