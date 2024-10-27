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

type RegisterController struct {
	jwtService httpinterfaces.JWTService
	interactor *interactors.RegisterInteractor
}

func NewRegisterController(jwtService httpinterfaces.JWTService, interactor *interactors.RegisterInteractor) *RegisterController {
	return &RegisterController{
		jwtService: jwtService,
		interactor: interactor,
	}
}

func (c *RegisterController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		case errors.Is(err, domainmodels.ErrLoginTaken):
			http.Error(w, "Login taken", http.StatusConflict)
		case errors.Is(err, domainmodels.ErrLoginValidation):
			http.Error(w, "Login invalid", http.StatusBadRequest)
		case errors.Is(err, domainmodels.ErrPasswordComplexity):
			http.Error(w, "Password complexity not achieved", http.StatusBadRequest)
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
