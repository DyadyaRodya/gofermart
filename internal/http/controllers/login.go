package controllers

import (
	"encoding/json"
	"errors"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	httpdto "github.com/DyadyaRodya/gofermart/internal/http/dto"
	httpinterfaces "github.com/DyadyaRodya/gofermart/internal/http/interfaces"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	"net/http"
)

type LoginController struct {
	jwtService httpinterfaces.JWTService
	interactor *interactors.LoginInteractor
}

func NewLoginController(jwtService httpinterfaces.JWTService, interactor *interactors.LoginInteractor) *LoginController {
	return &LoginController{
		jwtService: jwtService,
		interactor: interactor,
	}
}

func (c *LoginController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var credsDTO httpdto.CredentialsDTO
	err := json.NewDecoder(r.Body).Decode(&credsDTO)
	if err != nil {
		http.Error(w, "You need to provide login and password in JSON body", http.StatusBadRequest)
		return
	}
	creds := credsDTO.ConvertToInteractorsDTO()
	userInfo, err := c.interactor.Handle(ctx, creds)

	if err != nil {
		switch {
		case errors.Is(err, domainmodels.ErrWrongCredentials):
			http.Error(w, "Wrong login or password", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := c.jwtService.NewUserToken(userInfo, httpdto.TTL)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	http.SetCookie(w, httpdto.NewAuthCookie(token, httpdto.CookieName, httpdto.Path, httpdto.TTL))
	w.WriteHeader(http.StatusOK)
}
