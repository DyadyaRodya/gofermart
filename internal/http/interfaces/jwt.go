package interfaces

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"time"
)

type JWTService interface {
	NewUserToken(userInfo *domainmodels.UserInfo, ttl time.Duration) (string, error)
	ParseToken(tokenString string) (userUUID, login string)
}
