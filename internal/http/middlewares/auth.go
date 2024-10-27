package middlewares

import (
	"context"
	"errors"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	httpdto "github.com/DyadyaRodya/gofermart/internal/http/dto"
	httpinterfaces "github.com/DyadyaRodya/gofermart/internal/http/interfaces"
	"github.com/DyadyaRodya/gofermart/internal/interactors/interfaces"
	"net/http"
)

type AuthMiddleware struct {
	jwtService httpinterfaces.JWTService
	repo       interfaces.Repository
}

func NewAuthMiddleware(jwtService httpinterfaces.JWTService, repo interfaces.Repository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		repo:       repo,
	}
}

func (m *AuthMiddleware) WithAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(httpdto.CookieName)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		var userUUID, login string

		if err == nil {
			userUUID, login = m.jwtService.ParseToken(token.Value)
		}

		if userUUID == "" || login == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		dbSess, err := m.repo.NewSession(ctx)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer dbSess.Close(ctx)
		userInfo, err := dbSess.GetUserByLogin(ctx, login)
		if err != nil && !errors.Is(err, domainmodels.ErrUserNotFound) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if userInfo == nil || errors.Is(err, domainmodels.ErrUserNotFound) || userInfo.UUID != userUUID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		newCTX := context.WithValue(ctx, domainmodels.UserInfo{}, userInfo)
		newToken, err := m.jwtService.NewUserToken(userInfo, httpdto.TTL)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		http.SetCookie(w, httpdto.NewAuthCookie(newToken, httpdto.CookieName, httpdto.Path, httpdto.TTL))

		next.ServeHTTP(w, r.WithContext(newCTX))
	}
}
