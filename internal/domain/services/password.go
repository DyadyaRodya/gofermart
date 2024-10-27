package services

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"github.com/DyadyaRodya/gofermart/internal/domain/models"
	"unicode"
)

type PasswordComplexityConfig struct {
	Length          int
	NumberOfDigits  int
	NumberOfUpper   int
	NumberOfLower   int
	NumberOfSpecial int
}

type PasswordDomainService struct {
	config *PasswordComplexityConfig
}

func NewPasswordDomainService(config *PasswordComplexityConfig) *PasswordDomainService {
	return &PasswordDomainService{
		config: config,
	}
}

func (p *PasswordDomainService) Hash(password, salt string) (string, error) {
	saltBytes := []byte(salt)
	var passwordBytes = []byte(password)
	passwordBytes = append(passwordBytes, saltBytes...)

	var sha512Hasher = sha512.New()
	_, err := sha512Hasher.Write(passwordBytes)
	if err != nil {
		return "", errors.Join(err, models.ErrPasswordHashGeneration)
	}
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex, nil
}

func (p *PasswordDomainService) Compare(currPassword, hashedPassword, salt string) (bool, error) {
	currPasswordHash, err := p.Hash(currPassword, salt)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(currPasswordHash)) == 1, nil
}

func (p *PasswordDomainService) Validate(password string) bool {
	numberOfDigits := 0
	numberOfUpper := 0
	numberOfLower := 0
	numberOfSpecial := 0

	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			numberOfDigits++
		case unicode.IsUpper(char):
			numberOfUpper++
		case unicode.IsLower(char):
			numberOfLower++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			numberOfSpecial++
		default:
			return false // unsupported symbol
		}
	}

	return len(password) >= p.config.Length &&
		numberOfDigits >= p.config.NumberOfDigits &&
		numberOfUpper >= p.config.NumberOfUpper &&
		numberOfLower >= p.config.NumberOfLower &&
		numberOfSpecial >= p.config.NumberOfSpecial
}

func (p *PasswordDomainService) GenerateRandomSalt(saltSize int) (string, error) {
	var saltBytes = make([]byte, saltSize)
	_, err := rand.Read(saltBytes[:])
	if err != nil {
		return "", errors.Join(err, models.ErrSaltGeneration)
	}
	return string(saltBytes), nil

}
