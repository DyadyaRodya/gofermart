package services

import (
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type (
	Claims struct {
		jwt.RegisteredClaims
		UserUUID string
		Login    string
	}
	JWTDomainService struct {
		secretKey []byte
	}
)

func NewJWTDomainService(secretKey []byte) *JWTDomainService {
	return &JWTDomainService{
		secretKey: secretKey,
	}
}

func (j *JWTDomainService) NewUserToken(userInfo *domainmodels.UserInfo, ttl time.Duration) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserUUID:         userInfo.UUID,
		Login:            userInfo.Login,
	}

	if ttl > 0 { // allow to be infinite if ttl == 0
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (j *JWTDomainService) ParseToken(tokenString string) (userUUID, login string) {
	userUUID, login = "", ""

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return
	}

	if !token.Valid {
		return
	}
	userUUID, login = claims.UserUUID, claims.Login
	return
}
